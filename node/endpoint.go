////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server gRPC endpoints

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handle a Broadcasted Ask Online event
func (s *Comms) AskOnline(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	pingMsg := &pb.Ping{}
	err := ptypes.UnmarshalAny(msg.Message, pingMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &pb.Ack{}, s.handler.AskOnline(pingMsg, authMsg)
}

// Handles validation of reverse-authentication tokens
func (s *Comms) AuthenticateToken(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	return &pb.Ack{}, s.ValidateToken(msg)
}

// Handles reception of reverse-authentication token requests
func (s *Comms) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	token, err := s.GenerateToken()
	return &pb.AssignToken{
		Token: token,
	}, err
}

// Handle a NewRound event
func (s *Comms) CreateNewRound(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {

	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	roundInfoMsg := &pb.RoundInfo{}
	err := ptypes.UnmarshalAny(msg.Message, roundInfoMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Call the server handler to start a new round
	return &pb.Ack{}, s.handler.CreateNewRound(roundInfoMsg, authMsg)
}

// PostNewBatch polls the first node and sends a batch when it is ready
func (s *Comms) PostNewBatch(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	batchMsg := &pb.Batch{}
	err := ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Call the server handler to post a new batch
	err = s.handler.PostNewBatch(batchMsg, authMsg)

	return &pb.Ack{}, err
}

// Handle a Phase event
func (s *Comms) PostPhase(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack,
	error) {
	// Call the server handler with the msg
	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	batchMsg := &pb.Batch{}
	err := ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	s.handler.PostPhase(batchMsg, authMsg)
	return &pb.Ack{}, nil
}

// Handle a phase event using a stream server
func (s *Comms) StreamPostPhase(server pb.Node_StreamPostPhaseServer) error {
	return s.handler.StreamPostPhase(server)
}

// Handle a PostRoundPublicKey message
func (s *Comms) PostRoundPublicKey(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	// Call the server handler that receives the key share
	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	publicKeyMsg := &pb.RoundPublicKey{}
	err := ptypes.UnmarshalAny(msg.Message, publicKeyMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	s.handler.PostRoundPublicKey(publicKeyMsg, authMsg)
	return &pb.Ack{}, nil
}

// GetBufferInfo returns buffer size (number of completed precomputations)
func (s *Comms) GetRoundBufferInfo(ctx context.Context,
	msg *pb.AuthenticatedMessage) (
	*pb.RoundBufferInfo, error) {

	authMsg := s.AuthenticatedReceiver(msg)
	bufSize, err := s.handler.GetRoundBufferInfo(authMsg)
	if bufSize < 0 {
		bufSize = 0
	}
	size := uint32(bufSize)
	return &pb.RoundBufferInfo{RoundBufferSize: size}, err
}

// Handles Registration Nonce Communication
func (s *Comms) RequestNonce(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Nonce, error) {
	//Create an auth object
	authMsg := s.AuthenticatedReceiver(msg)

	//Marshall the any message to the message type needed
	nonceRequest := &pb.NonceRequest{}
	err := ptypes.UnmarshalAny(msg.Message, nonceRequest)
	if err != nil {
		return nil, err
	}

	// Obtain the nonce by passing to server
	nonce, pk, err := s.handler.RequestNonce(nonceRequest.GetSalt(),
		nonceRequest.GetClientRSAPubKey(), nonceRequest.GetClientDHPubKey(),
		nonceRequest.GetClientSignedByServer().Signature,
		nonceRequest.GetRequestSignature().Signature, authMsg)

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the NonceMessage
	return &pb.Nonce{
		Nonce:    nonce,
		DHPubKey: pk,
		Error:    errMsg,
	}, err
}

// Handles Registration Nonce Confirmation
func (s *Comms) ConfirmRegistration(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.RegistrationConfirmation, error) {
	authMsg := s.AuthenticatedReceiver(msg)

	//Marshall the any message to the message type needed
	regConfirmRequest := &pb.RequestRegistrationConfirmation{}
	err := ptypes.UnmarshalAny(msg.Message, regConfirmRequest)
	if err != nil {
		return nil, err
	}

	// Obtain signed client public key by passing to server
	signature, err := s.handler.ConfirmRegistration(regConfirmRequest.GetUserID(),
		regConfirmRequest.NonceSignedByClient.Signature, authMsg)

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the RegistrationConfirmation
	return &pb.RegistrationConfirmation{
		ClientSignedByServer: &pb.RSASignature{
			Signature: signature,
		},
		Error: errMsg,
	}, err
}

// PostPrecompResult sends final Message and AD precomputations.
func (s *Comms) PostPrecompResult(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Ack, error) {

	authMsg := s.AuthenticatedReceiver(msg)

	//Marshall the any message to the message type needed
	batchMsg := &pb.Batch{}
	err := ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, err
	}

	// Call the server handler to start a new round
	err = s.handler.PostPrecompResult(batchMsg.GetRound().GetID(),
		batchMsg.Slots, authMsg)
	return &pb.Ack{}, err
}

// FinishRealtime broadcasts to all nodes when the realtime is completed
func (s *Comms) FinishRealtime(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	// Call the server handler to finish realtime
	authMsg := s.AuthenticatedReceiver(msg)

	//Marshall the any message to the message type needed
	roundInfoMsg := &pb.RoundInfo{}
	err := ptypes.UnmarshalAny(msg.Message, roundInfoMsg)
	if err != nil {
		return nil, err
	}

	err = s.handler.FinishRealtime(roundInfoMsg, authMsg)

	return &pb.Ack{}, err
}

// GetCompletedBatch should return a completed batch that the calling gateway
// hasn't gotten before
func (s *Comms) GetCompletedBatch(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Batch, error) {

	authMsg := s.AuthenticatedReceiver(msg)
	return s.handler.GetCompletedBatch(authMsg)
}

func (s *Comms) GetMeasure(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.RoundMetrics, error) {
	authMsg := s.AuthenticatedReceiver(msg)

	//Marshall the any message to the message type needed
	roundInfoMsg := &pb.RoundInfo{}
	err := ptypes.UnmarshalAny(msg.Message, roundInfoMsg)
	if err != nil {
		return nil, err
	}

	rm, err := s.handler.GetMeasure(roundInfoMsg, authMsg)
	return rm, err
}

func (s *Comms) PollNdf(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.GatewayNdf, error) {
	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	pingMsg := &pb.Ping{}
	err := ptypes.UnmarshalAny(msg.Message, pingMsg)
	if err != nil {
		return nil, err
	}

	rm, err := s.handler.PollNdf(pingMsg, authMsg)
	return rm, err
}

func (s *Comms) SendRoundTripPing(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	authMsg := s.AuthenticatedReceiver(msg)
	//Marshall the any message to the message type needed
	roundTripPing := &pb.RoundTripPing{}
	err := ptypes.UnmarshalAny(msg.Message, roundTripPing)
	if err != nil {
		return nil, err
	}

	err = s.handler.SendRoundTripPing(roundTripPing, authMsg)
	return &pb.Ack{}, err
}
