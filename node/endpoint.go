////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server GRPC endpoints

package node

// TODO: A lot of message types from gRPC are passed through, and a number of
//       errors that can occur are not accounted for.

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handle a Broadcasted Ask Online event
func (s *server) AskOnline(ctx context.Context, msg *pb.Ping) (
	*pb.Ack, error) {
	return &pb.Ack{}, nil
}

// Handle a Roundtrip ping event
func (s *server) RoundtripPing(ctx context.Context, msg *pb.TimePing) (
	*pb.Ack, error) {
	serverHandler.RoundtripPing(msg)
	return &pb.Ack{}, nil
}

// Handle a broadcasted ServerMetric event
func (s *server) ServerMetrics(ctx context.Context, msg *pb.ServerMetricsMessage) (
	*pb.Ack, error) {
	serverHandler.ServerMetrics(msg)
	return &pb.Ack{}, nil
}

// Handle a NewRound event
func (s *server) NewRound(ctx context.Context,
	msg *pb.CmixBatch) (*pb.Ack, error) {
	// Call the server handler to start a new round
	serverHandler.NewRound(msg.RoundID)
	return &pb.Ack{}, nil
}

// Handle a Phase event
func (s *server) Phase(ctx context.Context, msg *pb.CmixBatch) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.Phase(msg)
	return &pb.Ack{}, nil
}

// Handle a StartRound event
func (s *server) StartRound(ctx context.Context,
	msg *pb.InputMessages) (*pb.Ack, error) {
	serverHandler.StartRound(msg)
	return &pb.Ack{}, nil
}

// GetBufferInfo returns buffer size (number of completed precomputations)
func (s *server) GetRoundBufferInfo(ctx context.Context, msg *pb.Ping) (
	*pb.RoundBufferInfo, error) {
	bufSize, err := serverHandler.GetRoundBufferInfo()
	if bufSize < 0 {
		bufSize = 0
	}
	size := uint32(bufSize)
	return &pb.RoundBufferInfo{RoundBufferSize: size}, err
}

// Handles Registration Nonce Communication
func (s *server) RequestNonce(ctx context.Context,
	msg *pb.RequestNonceMessage) (*pb.NonceMessage, error) {

	// Obtain the nonce by passing to server
	nonce, err := serverHandler.RequestNonce(msg.GetSalt(),
		msg.GetY(), msg.GetP(), msg.GetQ(),
		msg.GetG(), msg.GetHash(), msg.GetR(), msg.GetS())

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the NonceMessage
	return &pb.NonceMessage{
		Nonce: nonce,
		Error: errMsg,
	}, err
}

// Handles Registration Nonce Confirmation
func (s *server) ConfirmNonce(ctx context.Context,
	msg *pb.ConfirmNonceMessage) (*pb.RegistrationConfirmation, error) {

	// Obtain signed client public key by passing to server
	hash, R, S, Y, P, Q, G, err := serverHandler.ConfirmNonce(msg.GetHash(),
		msg.GetR(), msg.GetS())

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the RegistrationConfirmation
	return &pb.RegistrationConfirmation{
		Hash:  hash,
		R:     R,
		S:     S,
		Y:     Y,
		P:     P,
		Q:     Q,
		G:     G,
		Error: errMsg,
	}, err
}
