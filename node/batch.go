////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains logic for batch-related comms

package node

import (
	"context"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

// ------------------------- PrecompTestBatchBroadcast Logic ---------------------------------------- //

// StreamPrecompTestBatch is a server to server broadcast. It simulates
// sending the completed batch of PrecompTestBatch, testing for connectivity.
func (s *Comms) StreamPrecompTestBatch(host *connect.Host, info *pb.RoundInfo,
	mockBatch *pb.CompletedBatch) error {
	// Retrieve the streaming service
	streamingClient, cancel, err := s.getPrecompTestBatchStreamClient(host, info)
	if err != nil {
		return errors.Errorf("Could not retrieve steaming service: %v",
			err)
	}
	defer cancel()

	// Stream each slot
	for i, slot := range mockBatch.Slots {
		if err = streamingClient.Send(slot); err != nil {
			if err == io.EOF {
				// Attempt to read an error
				eofAck, eofErr := streamingClient.CloseAndRecv()
				if eofErr != nil {
					err = errors.Wrap(err, eofErr.Error())
				} else {
					err = errors.Wrap(err, eofAck.Error)
				}
			}
			return errors.Errorf("Could not stream slot (%d/%d) "+
				"for round %d: %v", i, len(mockBatch.Slots), info.ID, err)
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

// getPrecompTestBatchStreamClient gets the streaming client
// using a header and returns the stream and the cancel context
// if there are no connection errors.
func (s *Comms) getPrecompTestBatchStreamClient(host *connect.Host,
	info *pb.RoundInfo) (pb.Node_PrecompTestBatchClient,
	context.CancelFunc, error) {

	ctx, cancel := s.getPrecompTestBatchContext(info)

	streamClient, err := s.getPrecompTestBatchStream(host, ctx)
	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil

}

// getPrecompTestBatchContext is given roundInfo as a header,
// and creates a streaming context. It adds the header to the context
// and returns the context with the header and a cancel func.
func (s *Comms) getPrecompTestBatchContext(info *pb.RoundInfo) (context.Context,
	context.CancelFunc) {
	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	encodedStr := base64.StdEncoding.EncodeToString([]byte(info.String()))

	// Add batch information to streaming context
	ctx = metadata.AppendToOutgoingContext(ctx, pb.PrecompTestBatchHeader,
		encodedStr)

	return ctx, cancel
}

// getPrecompTestBatchStream returns the streaming client for PrecompTestBatch.
func (s *Comms) getPrecompTestBatchStream(host *connect.Host,
	ctx context.Context) (pb.Node_PrecompTestBatchClient, error) {

	// Create the Stream Function
	f := func(conn *grpc.ClientConn) (interface{}, error) {

		// Add authentication information to streaming context
		ctx = s.PackAuthenticatedContext(host, ctx)

		// Get the stream client
		streamClient, err := pb.NewNodeClient(conn).PrecompTestBatch(ctx)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return streamClient, nil
	}

	jww.TRACE.Printf("Streaming PrecompTestBatch")

	// Execute the Stream function
	resultClient, err := s.ProtoComms.Stream(host, f)
	if err != nil {
		return nil, err
	}

	result := resultClient.(pb.Node_PrecompTestBatchClient)
	return result, nil

}

// PrecompTestBatch is the reception handler for StreamPrecompTestBatch.
func (s *Comms) PrecompTestBatch(stream pb.Node_PrecompTestBatchServer) error {
	// Extract the authentication info
	authMsg, err := connect.UnpackAuthenticatedContext(stream.Context())
	if err != nil {
		return errors.Errorf("Unable to extract authentication info: %+v", err)
	}

	authState, err := s.AuthenticatedReceiver(authMsg, stream.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	info, err := GetPrecompTestBatchStreamHeader(stream)
	if err != nil {
		return errors.WithMessage(err, "Could not get test batch stream header")
	}

	return s.handler.PrecompTestBatch(stream, info, authState)
}

// GetPrecompTestBatchStreamHeader gets the header in the metadata from
// the server stream and converts it to a mixmessages.RoundInfo message.
func GetPrecompTestBatchStreamHeader(stream pb.Node_PrecompTestBatchServer) (*pb.RoundInfo, error) {
	// Obtain the headers from server metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	// Unmarshall the header into a message
	marshledBatch, err := base64.StdEncoding.DecodeString(md.Get(pb.PrecompTestBatchHeader)[0])
	if err != nil {
		return nil, err
	}
	batchInfo := &pb.RoundInfo{}
	err = proto.UnmarshalText(string(marshledBatch), batchInfo)
	if err != nil {
		return nil, err
	}

	return batchInfo, nil
}

// ------------------------- UploadUnmixedBatch Logic ---------------------------------------- //

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
func (s *Comms) DownloadMixedBatch(authMsg *messages.AuthenticatedMessage,
	stream pb.Node_DownloadMixedBatchServer) error {

	authState, err := s.AuthenticatedReceiver(authMsg, stream.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	batchInfo := &pb.BatchReady{}
	err = ptypes.UnmarshalAny(authMsg.Message, batchInfo)
	if err != nil {
		return err
	}

	return s.handler.DownloadMixedBatch(stream, batchInfo, authState)
}
