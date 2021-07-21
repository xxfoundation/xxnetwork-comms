///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains logic for batch-related comms

package node

import (
	"context"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UploadUnmixedBatch is the handler for gateway sending a batch to its node
func (s *Comms) UploadUnmixedBatch(server pb.Node_UploadUnmixedBatchServer) error {
	// Extract the authentication info
	authMsg, err := connect.UnpackAuthenticatedContext(server.Context())
	if err != nil {
		return errors.Errorf("Unable to extract authentication info: %+v", err)
	}

	authState, err := s.AuthenticatedReceiver(authMsg, server.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Verify the message authentication
	return s.handler.UploadUnmixedBatch(server, authState)
}

// GetUnmixedBatchStreamHeader gets the header in the metadata from
// the server stream and returns it or an error if it fails.
func GetUnmixedBatchStreamHeader(stream pb.Node_UploadUnmixedBatchServer) (*pb.BatchInfo, error) {

	// Obtain the headers from server metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	// Unmarshall the header into a message
	marshledBatch, err := base64.StdEncoding.DecodeString(md.Get(pb.UnmixedBatchHeader)[0])
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


// ------------------------- DownloadMixedBatch Logic ---------------------------------------- //

// DownloadMixedBatch streams the slots in the completed batch to the gateway
func (s *Comms) DownloadMixedBatch(host *connect.Host,
	info pb.BatchInfo, batch *pb.CompletedBatch) error {

	// Retrieve the streaming service
	streamingClient, cancel, err := s.getMixedBatchStreamClient(
		host, info)
	if err != nil {
		return errors.Errorf("Could not retrieve steaming service: %v", err)
	}
	defer cancel()

	// Stream each slot
	for i, slot := range batch.Slots {
		if err = streamingClient.Send(slot); err != nil {
			return errors.Errorf("Could not stream "+
				"slot (%d/%d) for round %d: %v",
				i, len(batch.Slots), batch.RoundID, err)
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
func (s *Comms) getMixedBatchStreamClient(host *connect.Host,
	header pb.BatchInfo) (pb.Gateway_DownloadMixedBatchClient, context.CancelFunc, error) {

	ctx, cancel := s.getMixedBatchStreamContext(&header)

	streamClient, err := s.getMixedBatchStream(host, ctx)
	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil
}

// getMixedBatchStreamContext is given batchInfo header
// and creates a streaming context, adds the header to the context
// and returns the context with the header and a cancel func
func (s *Comms) getMixedBatchStreamContext(batchInfo *pb.BatchInfo) (
	context.Context, context.CancelFunc) {

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	encodedStr := base64.StdEncoding.EncodeToString([]byte(batchInfo.String()))

	// Add batch information to streaming context
	ctx = metadata.AppendToOutgoingContext(ctx, pb.MixedBatchHeader, encodedStr)

	return ctx, cancel
}

// getMixedBatchStream uses an id and streaming context to retrieve
// a Gateway_DownloadMixedBatchClient object otherwise it returns
// an error if the connection is unavailable
func (s *Comms) getMixedBatchStream(host *connect.Host,
	ctx context.Context) (pb.Gateway_DownloadMixedBatchClient, error) {

	// Create the Stream Function
	f := func(conn *grpc.ClientConn) (interface{}, error) {

		// Add authentication information to streaming context
		ctx = s.PackAuthenticatedContext(host, ctx)

		// Send the message
		stream, err := pb.NewGatewayClient(conn).DownloadMixedBatch(ctx)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		return stream, nil
	}

	jww.TRACE.Printf("Streaming DownloadMixedBatch")

	// Execute the Stream function
	resultClient, err := s.ProtoComms.Stream(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := resultClient.(pb.Gateway_DownloadMixedBatchClient)
	return result, nil
}
