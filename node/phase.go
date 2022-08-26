///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains server -> server functionality for precomputation operations

package node

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
	"google.golang.org/grpc/metadata"
)

// Server -> Server Send Function
func (s *Comms) SendPostPhase(host *connect.Host,
	message *pb.Batch) (*messages.Ack, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Format to authenticated message type
		authMsg, err := s.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Send the message
		resultMsg, err := pb.NewNodeClient(conn.GetGrpcConn()).
			PostPhase(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Post Phase message: %+v", message)
	resultMsg, err := s.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// GetPostPhaseStreamClient gets the streaming client
// using a header and returns the stream and the cancel context
// if there are no connection errors
func (s *Comms) GetPostPhaseStreamClient(host *connect.Host,
	header pb.BatchInfo) (pb.Node_StreamPostPhaseClient, context.CancelFunc, error) {

	ctx, cancel := s.getPostPhaseStreamContext(&header)

	streamClient, err := s.getPostPhaseStream(host, ctx)
	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil
}

// getPostPhaseStreamContext is given batchInfo PostPhase header
// and creates a streaming context, adds the header to the context
// and returns the context with the header and a cancel func
func (s *Comms) getPostPhaseStreamContext(batchInfo *pb.BatchInfo) (
	context.Context, context.CancelFunc) {

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	encodedStr := base64.StdEncoding.EncodeToString([]byte(batchInfo.String()))

	// Add batch information to streaming context
	ctx = metadata.AppendToOutgoingContext(ctx, pb.PostPhaseHeader, encodedStr)

	return ctx, cancel
}

// getPostPhaseStream uses an id and streaming context to retrieve
// a Node_StreamPostPhaseClient object otherwise it returns
// an error if the connection is unavailable
func (s *Comms) getPostPhaseStream(host *connect.Host,
	ctx context.Context) (pb.Node_StreamPostPhaseClient, error) {

	// Create the Stream Function
	f := func(conn connect.Connection) (interface{}, error) {

		// Add authentication information to streaming context
		ctx = s.PackAuthenticatedContext(host, ctx)

		// Get the stream client
		streamClient, err := pb.NewNodeClient(conn.GetGrpcConn()).
			StreamPostPhase(ctx)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return streamClient, nil
	}

	// Execute the Stream function
	resultClient, err := s.Stream(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := resultClient.(pb.Node_StreamPostPhaseClient)
	return result, nil
}

// Gets the header in the metadata from the server stream
// and returns it or an error if it fails.
func GetPostPhaseStreamHeader(stream pb.Node_StreamPostPhaseServer) (*pb.BatchInfo, error) {

	// Obtain the headers from server metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	// Unmarshall the header into a message

	marshledBatch, err := base64.StdEncoding.DecodeString(md.Get(pb.PostPhaseHeader)[0])
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
