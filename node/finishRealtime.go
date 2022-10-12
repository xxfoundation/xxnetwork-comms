////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// This file contains logic for streaming the slots of a completed batch to
// all other nodes in the round.

package node

import (
	"context"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

/* ------------------- Broadcast functions ------------------- */

// SendFinishRealtime is a node to node comm which streams a completed batch
// to all nodes within a round.
func (s *Comms) SendFinishRealtime(host *connect.Host,
	roundInfo *pb.RoundInfo, batch *pb.CompletedBatch) (*messages.Ack, error) {

	// Retrieve the streaming service
	streamingClient, cancel, err := s.getFinishRealtimeStreamClient(
		host, roundInfo)
	if err != nil {
		return nil, errors.Errorf("Could not retrieve steaming service: %v",
			err)
	}
	defer cancel()

	// Stream each slot
	for i, slot := range batch.Slots {
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
			return nil, errors.Errorf("Could not stream slot (%d/%d) to %s"+
				"for round %d to %s: %v", i, len(batch.Slots), host.GetId(),
				roundInfo.ID, host.GetId(), err)
		}
	}

	// Receive ack and cancel client streaming context
	ack, err := streamingClient.CloseAndRecv()
	if err != nil {
		return nil, errors.Errorf("Could not receive final "+
			"acknowledgement on streaming batch: %v", err)
	}

	if ack != nil && ack.Error != "" {
		return nil, errors.Errorf("Remote Server Error: %v", ack.Error)
	}

	return nil, nil
}

// getFinishRealtimeStreamClient gets the streaming client
// using a header and returns the stream and the cancel context
// if there are no connection errors.
func (s *Comms) getFinishRealtimeStreamClient(host *connect.Host,
	info *pb.RoundInfo) (pb.Node_FinishRealtimeClient, context.CancelFunc,
	error) {

	ctx, cancel := s.getFinishRealtimeContext(info)

	streamClient, err := s.getFinishRealtimeStream(host, ctx)
	if err != nil {
		return nil, nil, err
	}

	return streamClient, cancel, nil
}

// getFinishRealtimeContext is given roundInfo as a header,
// and creates a streaming context. It adds the header to the context
// and returns the context with the header and a cancel func.
func (s *Comms) getFinishRealtimeContext(info *pb.RoundInfo) (context.Context,
	context.CancelFunc) {

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	encodedStr := base64.StdEncoding.EncodeToString([]byte(info.String()))

	// Add batch information to streaming context
	ctx = metadata.AppendToOutgoingContext(ctx, pb.FinishRealtimeHeader,
		encodedStr)

	return ctx, cancel
}

// getFinishRealtimeStream returns the streaming client for FinishRealtime.
func (s *Comms) getFinishRealtimeStream(host *connect.Host,
	ctx context.Context) (pb.Node_FinishRealtimeClient, error) {

	// Create the Stream Function
	f := func(conn *grpc.ClientConn) (interface{}, error) {

		// Add authentication information to streaming context
		ctx = s.PackAuthenticatedContext(host, ctx)

		// Get the stream client
		streamClient, err := pb.NewNodeClient(conn).FinishRealtime(ctx)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return streamClient, nil
	}

	jww.TRACE.Printf("Streaming FinishRealtime")

	// Execute the Stream function
	resultClient, err := s.ProtoComms.Stream(host, f)
	if err != nil {
		return nil, err
	}

	result := resultClient.(pb.Node_FinishRealtimeClient)
	return result, nil

}

/* ------------------- Receive functions ------------------- */

// FinishRealtime broadcasts to all nodes when the realtime is completed
func (s *Comms) FinishRealtime(stream pb.Node_FinishRealtimeServer) error {
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
	info, err := GetFinishRealtimeStreamHeader(stream)
	if err != nil {
		return errors.WithMessage(err, "Could not get realtime stream header")
	}

	// Handle
	return s.handler.FinishRealtime(info, stream, authState)
}

// GetFinishRealtimeStreamHeader gets the header in the metadata from
// the server stream and converts it to a mixmessages.RoundInfo message.
func GetFinishRealtimeStreamHeader(stream pb.Node_FinishRealtimeServer) (*pb.RoundInfo, error) {

	// Obtain the headers from server metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	// Unmarshall the header into a message
	marshledBatch, err := base64.StdEncoding.DecodeString(md.Get(pb.FinishRealtimeHeader)[0])
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
