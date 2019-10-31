////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> server functionality for precomputation operations

package node

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc/metadata"
)

// Server -> Server Send Function
func (s *NodeComms) SendPostPhase(connInfo *connect.Host,
	message *pb.Batch) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).PostPhase(ctx, message,
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// GetPostPhaseStreamClient gets the streaming client
// using a header and returns the stream and the cancel context
// if there are no connection errors
func (s *NodeComms) GetPostPhaseStreamClient(connInfo *connect.Host,
	header pb.BatchInfo) (pb.Node_StreamPostPhaseClient, context.CancelFunc, error) {

	ctx, cancel := s.getPostPhaseStreamContext(header)
	streamClient, err := s.getPostPhaseStream(connInfo, ctx)
	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil
}

// getPostPhaseStreamContext is given batchInfo PostPhase header
// and creates a streaming context, adds the header to the context
// and returns the context with the header and a cancel func
func (s *NodeComms) getPostPhaseStreamContext(batchInfo pb.BatchInfo) (
	context.Context, context.CancelFunc) {

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	// Create a new context with some metadata
	// using the batch info batchInfo
	ctx = metadata.AppendToOutgoingContext(ctx, "batchinfo", batchInfo.String())

	return ctx, cancel
}

// getPostPhaseStream uses an id and streaming context to retrieve
// a Node_StreamPostPhaseClient object otherwise it returns
// an error if the connection is unavailable
func (s *NodeComms) getPostPhaseStream(connInfo *connect.Host,
	ctx context.Context) (pb.Node_StreamPostPhaseClient, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Get the stream client using streaming context
	streamClient, err := pb.NewNodeClient(conn.Connection).StreamPostPhase(ctx,
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return streamClient, nil
}

// GetPostPhaseStreamHeader gets the header
// in the metadata from the server stream
// and returns it with an error if it fails.
func GetPostPhaseStreamHeader(stream pb.Node_StreamPostPhaseServer) (*pb.BatchInfo, error) {

	// Unmarshal header into batch info
	batchInfo := pb.BatchInfo{}

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header %v")
	}

	err := proto.UnmarshalText(md.Get("batchinfo")[0], &batchInfo)
	if err != nil {
		return nil, err
	}

	return &batchInfo, nil
}
