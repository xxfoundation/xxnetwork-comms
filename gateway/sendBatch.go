////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server functionality

package gateway

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Gateway -> Server Send Function
func (g *GatewayComms) PostNewBatch(connInfo *connect.Host,
	messages *pb.Batch) error {

	// Obtain the connection
	conn, err := g.ObtainConnection(connInfo)
	if err != nil {
		return err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	_, err = pb.NewNodeClient(conn.Connection).PostNewBatch(ctx, messages)
	if err != nil {
		err = errors.New(err.Error())
	}

	return err
}

// GetRoundBufferInfo Asks the server for round buffer info, specifically how
// many rounds have gone through precomputation.
// Note that this function should block if the buffer size is 0
// This allows the caller to continuously poll without spinning too much.
func (g *GatewayComms) GetRoundBufferInfo(
	connInfo *connect.Host) (int, error) {

	// Initialize bufSize
	bufSize := 0

	// Obtain the connection
	conn, err := g.ObtainConnection(connInfo)
	if err != nil {
		return bufSize, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	bufInfo, err := pb.NewNodeClient(
		conn.Connection).GetRoundBufferInfo(ctx, &pb.RoundBufferInfo{})
	if err != nil {
		err = errors.New(err.Error())
	} else {
		bufSize = int(bufInfo.RoundBufferSize)
	}

	return bufSize, err
}

// Gateway -> Server Send Function
func (g *GatewayComms) GetCompletedBatch(
	connInfo *connect.Host) (*pb.Batch, error) {

	// Obtain the connection
	conn, err := g.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	batch, err := pb.NewNodeClient(conn.Connection).GetCompletedBatch(ctx, &pb.Ping{})
	if err != nil {
		err = errors.New(err.Error())
	}

	return batch, err
}
