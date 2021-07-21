///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains batch streaming functionality

package gateway

import (
	"context"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// --------------------------- UploadMixedBatch Logic ----------------------------------------//

// UploadUnmixedBatch streams the slots in the batch to the node
func (g *Comms) UploadUnmixedBatch(host *connect.Host,
	batchInfo pb.BatchInfo, batch *pb.Batch) error {
	// Retrieve the streaming service
	streamingClient, cancel, err := g.getUnmixedBatchStreamClient(
		host, batchInfo)
	if err != nil {
		return errors.Errorf("Could not retrieve steaming service: %v", err)
	}
	defer cancel()

	// Stream each slot
	for i, slot := range batch.Slots {
		if err = streamingClient.Send(slot); err != nil {
			return errors.Errorf("Could not stream "+
				"slot (%d/%d) for round %d: %v",
				i, len(batch.Slots), batch.Round.ID, err)
		}
	}

	// Receive ack and cancel client streaming context
	ack, err := streamingClient.CloseAndRecv()
	if err != nil {
		return errors.Errorf("Could not receive final "+
			"acknowledgement on streaming batch: %v", err)
	}

	if ack != nil && ack.Error != "" {
		return errors.Errorf("Remote Server Error: %v", ack.Error)
	}

	return nil
}

// getUnmixedBatchStreamClient gets the streaming client
// using a header and returns the stream and the cancel context
// if there are no connection errors
func (g *Comms) getUnmixedBatchStreamClient(host *connect.Host,
	header pb.BatchInfo) (pb.Node_UploadUnmixedBatchClient, context.CancelFunc, error) {

	ctx, cancel := g.getUnmixedBatchStreamContext(&header)

	streamClient, err := g.getUnmixedBatchStream(host, ctx)
	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil
}

// getUnmixedBatchStreamContext is given batchInfo header
// and creates a streaming context, adds the header to the context
// and returns the context with the header and a cancel func
func (g *Comms) getUnmixedBatchStreamContext(batchInfo *pb.BatchInfo) (
	context.Context, context.CancelFunc) {

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	encodedStr := base64.StdEncoding.EncodeToString([]byte(batchInfo.String()))

	// Add batch information to streaming context
	ctx = metadata.AppendToOutgoingContext(ctx, pb.UnmixedBatchHeader, encodedStr)

	return ctx, cancel
}

// getUnmixedBatchStream uses an id and streaming context to retrieve
// a Node_UploadUnmixedBatchClient object otherwise it returns
// an error if the connection is unavailable
func (g *Comms) getUnmixedBatchStream(host *connect.Host,
	ctx context.Context) (pb.Node_UploadUnmixedBatchClient, error) {

	// Create the Stream Function
	f := func(conn *grpc.ClientConn) (interface{}, error) {

		// Add authentication information to streaming context
		ctx = g.PackAuthenticatedContext(host, ctx)

		// Get the stream client
		streamClient, err := pb.NewNodeClient(conn).UploadUnmixedBatch(ctx)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return streamClient, nil
	}

	jww.TRACE.Printf("Streaming UploadUnmixedBatch")

	// Execute the Stream function
	resultClient, err := g.ProtoComms.Stream(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := resultClient.(pb.Node_UploadUnmixedBatchClient)
	return result, nil
}

// ------------------------- DownloadMixedBatch Logic ----------------------------------------//

// StartDownloadMixedBatch sends a request for streaming a completed batch.
func (g *Comms) StartDownloadMixedBatch(host *connect.Host, ready *pb.BatchReady) error {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := g.PackAuthenticatedMessage(ready, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).StartDownloadMixedBatch(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Poll message...")
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return err
	}

	// Marshall the result
	result := &messages.Ack{}
	return ptypes.UnmarshalAny(resultMsg, result)
}

// DownloadMixedBatch is the handler for server sending a completed batch to its gateway
func (g *Comms) DownloadMixedBatch(server pb.Gateway_DownloadMixedBatchServer) error {
	// Extract the authentication info
	authMsg, err := connect.UnpackAuthenticatedContext(server.Context())
	if err != nil {
		return errors.Errorf("Unable to extract authentication info: %+v", err)
	}

	authState, err := g.AuthenticatedReceiver(authMsg, server.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Verify the message authentication
	return g.handler.DownloadMixedBatch(server, authState)
}

// GetMixedBatchStreamHeader gets the header in the metadata from
// the server stream and returns it or an error if it fails.
func GetMixedBatchStreamHeader(stream pb.Gateway_DownloadMixedBatchServer) (*pb.BatchInfo, error) {
	// Obtain the headers from server metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	// Unmarshall the header into a message
	marshledBatch, err := base64.StdEncoding.DecodeString(md.Get(pb.MixedBatchHeader)[0])
	if err != nil {
		return nil, err
	}
	batchInfo := &pb.BatchInfo{}
	err = proto.UnmarshalText(string(marshledBatch), batchInfo)
	if err != nil {
		return nil, err
	}

	return batchInfo, nil
}
