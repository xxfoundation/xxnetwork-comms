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
	msg *pb.InitRound) (*pb.Ack, error) {
	// Call the server handler to start a new round
	serverHandler.NewRound(msg.RoundID)
	return &pb.Ack{}, nil
}

// Handle a PrecompDecrypt event
func (s *server) PrecompDecrypt(ctx context.Context,
	msg *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompDecrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompEncrypt event
func (s *server) PrecompEncrypt(ctx context.Context,
	msg *pb.PrecompEncryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompEncrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompReveal event
func (s *server) PrecompReveal(ctx context.Context,
	msg *pb.PrecompRevealMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompReveal(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompPermute event
func (s *server) PrecompPermute(ctx context.Context,
	msg *pb.PrecompPermuteMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompPermute(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompShare event
func (s *server) PrecompShare(ctx context.Context,
	msg *pb.PrecompShareMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompShare(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompShareInit event
func (s *server) PrecompShareInit(ctx context.Context,
	msg *pb.PrecompShareInitMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompShareInit(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompShareCompare event
func (s *server) PrecompShareCompare(ctx context.Context,
	msg *pb.PrecompShareCompareMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompShareCompare(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompShareConfirm event
func (s *server) PrecompShareConfirm(ctx context.Context,
	msg *pb.PrecompShareConfirmMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompShareConfirm(msg)
	return &pb.Ack{}, nil
}

// Handle a RealtimeDecrypt event
func (s *server) RealtimeDecrypt(ctx context.Context,
	msg *pb.RealtimeDecryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.RealtimeDecrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a RealtimeDecrypt event
func (s *server) RealtimeEncrypt(ctx context.Context,
	msg *pb.RealtimeEncryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.RealtimeEncrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a RealtimePermute event
func (s *server) RealtimePermute(ctx context.Context,
	msg *pb.RealtimePermuteMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.RealtimePermute(msg)
	return &pb.Ack{}, nil
}

// Handle a SetPublicKey event
func (s *server) SetPublicKey(ctx context.Context,
	msg *pb.PublicKeyMessage) (*pb.Ack, error) {
	serverHandler.SetPublicKey(msg.RoundID, msg.PublicKey)
	return &pb.Ack{}, nil
}

// Handle a StartRound event
func (s *server) StartRound(ctx context.Context,
	msg *pb.InputMessages) (*pb.Ack, error) {
	serverHandler.StartRound(msg)
	return &pb.Ack{}, nil
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
	hash, R, S, err := serverHandler.ConfirmNonce(msg.GetHash(),
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
		Error: errMsg,
	}, err
}
