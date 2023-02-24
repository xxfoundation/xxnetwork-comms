////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"google.golang.org/grpc"
	"io"
	"strconv"
	"time"
)

// SendPutMessage Client -> Gateway Send Function
func (c *Comms) SendPutMessage(host *connect.Host, message *pb.GatewaySlot,
	timeout time.Duration) (*pb.GatewaySlotResponse, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContextWithTimeout(timeout)
		defer cancel()

		// Send the message
		var resultMsg = &pb.GatewaySlotResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.Gateway/PutMessage", message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				PutMessage(ctx, message)
		}

		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Put message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.GatewaySlotResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// SendPutManyMessages Client -> Gateway Send Function
func (c *Comms) SendPutManyMessages(host *connect.Host,
	messages *pb.GatewaySlots, timeout time.Duration) (
	*pb.GatewaySlotResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContextWithTimeout(timeout)
		defer cancel()

		// Send the message
		var resultMsg = &pb.GatewaySlotResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.Gateway/PutManyMessages",
				messages, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				PutManyMessages(ctx, messages)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending PutManyMessages: %+v", messages)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.GatewaySlotResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// SendRequestClientKeyMessage Client -> Gateway Send Function
func (c *Comms) SendRequestClientKeyMessage(host *connect.Host,
	message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &pb.SignedKeyResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.Gateway/RequestClientKey",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				RequestClientKey(ctx, message)
		}

		// Make sure there are no errors with sending the message
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Client Key message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.SignedKeyResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) BatchNodeRegistration(host *connect.Host,
	message *pb.SignedClientBatchKeyRequest) (*pb.SignedBatchKeyResponse, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &pb.SignedBatchKeyResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.Gateway/BatchNodeRegistration",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				BatchNodeRegistration(ctx, message)
		}

		// Make sure there are no errors with sending the message
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Client Key message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.SignedBatchKeyResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// SendPoll Client -> Gateway Send Function
// Returns a time.Time of the local clock (not netTime) when the comm was sent
// and a time.Duration representing the roundTripTime of the comm
func (c *Comms) SendPoll(host *connect.Host,
	message *pb.GatewayPoll) (*pb.GatewayPollResponse, time.Time, time.Duration, error) {
	// Set up the context with a timeout to ensure that streaming does not
	// block the follower
	ctx, cancel := connect.StreamingContextWithTimeout(10 * time.Second)
	defer cancel()

	var startTime time.Time

	// Create the Stream Function
	roundTripTime := time.Duration(0)
	f := func(conn connect.Connection) (interface{}, error) {
		// Send the message
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			clientStream, err := wc.NewServerStream(
				&grpc.StreamDesc{ServerStreams: true},
				"/mixmessages.Gateway/Poll")
			if err != nil {
				return nil, err
			}

			// use the local time NOT netTime because calculations for
			// clock skew will be done with this timestamp and including
			// the skew adjustment will break the calculation
			startTime = time.Now()
			err = clientStream.Send(ctx, message)
			roundTripTime = time.Now().Sub(startTime)
			if err != nil {
				return nil, err
			}
			return newServerStream(ctx, clientStream), nil
		} else {
			// use the local time NOT netTime because calculations for
			// clock skew will be done with this timestamp and including
			// the skew adjustment will break the calculation
			startTime = time.Now()
			clientStream, err := pb.NewGatewayClient(conn.GetGrpcConn()).
				Poll(ctx, message)
			roundTripTime = time.Now().Sub(startTime)
			if err != nil {
				return nil, err
			}
			return clientStream, nil
		}
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Poll message: %+v", message)
	resultClient, err := c.Stream(host, f)
	if err != nil {
		return nil, time.Time{}, 0, err
	}

	stream := resultClient.(pb.Gateway_PollClient)
	jww.DEBUG.Printf("Receiving chunks for gateway poll from %s", host.GetId().String())
	closeErr := stream.CloseSend()
	if closeErr != nil {
		return nil, time.Time{}, 0, wrapError(closeErr, "Unable to close send stream")
	}

	// Get the total number of chunks from the header
	md, err := stream.Header()
	if err != nil {
		closeErr = stream.RecvMsg(nil)
		return nil, time.Time{}, 0, wrapError(closeErr, "Could not "+
			"receive streaming header from %s: %s", host.GetId(), err)
	}

	// Check if metadata has the expected header
	chunkHeader := md.Get(pb.ChunkHeader)
	if len(chunkHeader) == 0 {
		closeErr = stream.RecvMsg(nil)
		return nil, time.Time{}, 0, wrapError(closeErr, pb.NoStreamingHeaderErr, host.GetId())
	}

	// Process header
	totalChunks, err := strconv.Atoi(chunkHeader[0])
	if err != nil {
		closeErr = stream.RecvMsg(nil)
		return nil, time.Time{}, 0, wrapError(closeErr, "Invalid header received: %v", err)
	}

	// Receive the chunks
	chunks := make([]*pb.StreamChunk, 0, totalChunks)
	chunk, err := stream.Recv()
	receivedChunks := 0
	for ; err == nil && receivedChunks <= totalChunks; chunk, err = stream.Recv() {
		chunks = append(chunks, chunk)
		receivedChunks++
	}
	if err != io.EOF { // EOF is an expected error after server-side has completed streaming
		return nil, time.Time{}, 0, errors.Errorf("Failed to "+
			"complete streaming, received %d of %d messages: %s",
			receivedChunks, totalChunks, err)
	}

	// Close stream once done
	closeErr = stream.RecvMsg(nil)
	if closeErr != io.EOF {
		return nil, time.Time{}, 0, errors.WithMessagef(closeErr, "Received error on "+
			"closing stream with %s", host.GetId())
	}

	// Assemble the result
	result := &pb.GatewayPollResponse{}
	return result, startTime, roundTripTime, pb.AssembleChunksIntoResponse(chunks, result)
}

// RequestHistoricalRounds Client -> Gateway Send Function
func (c *Comms) RequestHistoricalRounds(host *connect.Host,
	message *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &pb.HistoricalRoundsResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.Gateway/RequestHistoricalRounds",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				RequestHistoricalRounds(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Requesting Historical Rounds: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.HistoricalRoundsResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// RequestMessages Client -> Gateway Send Function
func (c *Comms) RequestMessages(host *connect.Host,
	message *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		var resultMsg = &pb.GetMessagesResponse{}
		var err error
		// Send the message
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.Gateway/RequestMessages", message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				RequestMessages(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Requesing Messages: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GetMessagesResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// RequestMessages Client -> Gateway Send Function
func (c *Comms) RequestBatchMessages(host *connect.Host,
	message *pb.GetMessagesBatch) (*pb.GetMessagesResponseBatch, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		var resultMsg = &pb.GetMessagesResponseBatch{}
		var err error
		// Send the message
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.Gateway/RequestBatchMessages", message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				RequestBatchMessages(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Requesting batch of Messages: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GetMessagesResponseBatch{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// GetGatewayTLSCertificate Client -> Gateway cert request
func (c *Comms) GetGatewayTLSCertificate(host *connect.Host,
	message *pb.RequestGatewayCert) (*pb.GatewayCertificate, error) {
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		var resultMsg = &pb.GatewayCertificate{}
		var err error
		// Send the message
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.Gateway/RequestTlsCert", message, resultMsg)
		} else {
			resultMsg, err = pb.NewGatewayClient(conn.GetGrpcConn()).
				RequestTlsCert(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Requesing TLS certificate from gateway: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GatewayCertificate{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

func wrapError(err error, s string, i ...interface{}) error {
	if err == nil {
		return errors.Errorf(s, i...)
	}
	return errors.Wrapf(err, s, i...)
}
