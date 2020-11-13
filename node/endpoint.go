///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains server gRPC endpoints

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"golang.org/x/net/context"
)

// Handle a Broadcasted Ask Online event
func (s *Comms) AskOnline(ctx context.Context, ping *messages.Ping) (*messages.Ack, error) {
	return &messages.Ack{}, s.handler.AskOnline()
}

// Handles validation of reverse-authentication tokens
func (s *Comms) AuthenticateToken(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	return &messages.Ack{}, s.ValidateToken(msg)
}

// Handles reception of reverse-authentication token requests
func (s *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := s.GenerateToken()
	return &messages.AssignToken{
		Token: token,
	}, err
}

// Handle a NewRound event
func (s *Comms) CreateNewRound(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unnmarshall the any message to the message type needed
	roundInfoMsg := &pb.RoundInfo{}
	err = ptypes.UnmarshalAny(msg.Message, roundInfoMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Call the server handler to start a new round
	return &messages.Ack{}, s.handler.CreateNewRound(roundInfoMsg, authState)
}

// PostNewBatch polls the first node and sends a batch when it is ready
func (s *Comms) PostNewBatch(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	// Unmarshall the any message to the message type needed
	batchMsg := &pb.Batch{}
	err = ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Call the server handler to post a new batch
	err = s.handler.PostNewBatch(batchMsg, authState)

	return &messages.Ack{}, err
}

// Handle a Phase event
func (s *Comms) PostPhase(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack,
	error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	// Unmarshall the any message to the message type needed
	batchMsg := &pb.Batch{}
	err = ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	// Call the server handler with the msg
	err = s.handler.PostPhase(batchMsg, authState)
	if err != nil {
		return &messages.Ack{}, err
	}
	return &messages.Ack{}, err
}

// Handle a phase event using a stream server
func (s *Comms) StreamPostPhase(server pb.Node_StreamPostPhaseServer) error {
	// Extract the authentication info
	authMsg, err := connect.UnpackAuthenticatedContext(server.Context())
	if err != nil {
		return errors.Errorf("Unable to extract authentication info: %+v", err)
	}

	authState, err := s.AuthenticatedReceiver(authMsg)
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Verify the message authentication
	return s.handler.StreamPostPhase(server, authState)
}

// Handle a PostRoundPublicKey message
func (s *Comms) PostRoundPublicKey(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	//Marshall the any message to the message type needed
	publicKeyMsg := &pb.RoundPublicKey{}
	err = ptypes.UnmarshalAny(msg.Message, publicKeyMsg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = s.handler.PostRoundPublicKey(publicKeyMsg, authState)
	if err != nil {
		return &messages.Ack{}, err
	}
	return &messages.Ack{}, err
}

// GetBufferInfo returns buffer size (number of completed precomputations)
func (s *Comms) GetRoundBufferInfo(ctx context.Context,
	msg *messages.AuthenticatedMessage) (
	*pb.RoundBufferInfo, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	bufSize, err := s.handler.GetRoundBufferInfo(authState)
	if bufSize < 0 {
		bufSize = 0
	}
	size := uint32(bufSize)
	return &pb.RoundBufferInfo{RoundBufferSize: size}, err
}

// Handles Registration Nonce Communication
func (s *Comms) RequestNonce(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*pb.Nonce, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Marshall the any message to the message type needed
	nonceRequest := &pb.NonceRequest{}
	err = ptypes.UnmarshalAny(msg.Message, nonceRequest)
	if err != nil {
		return nil, err
	}

	// Obtain the nonce by passing to server
	nonce, pk, err := s.handler.RequestNonce(nonceRequest.GetSalt(),
		nonceRequest.GetClientRSAPubKey(), nonceRequest.GetClientDHPubKey(),
		nonceRequest.GetClientSignedByServer().Signature,
		nonceRequest.GetRequestSignature().Signature, authState)

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
	msg *messages.AuthenticatedMessage) (*pb.RegistrationConfirmation, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	regConfirmRequest := &pb.RequestRegistrationConfirmation{}
	err = ptypes.UnmarshalAny(msg.Message, regConfirmRequest)
	if err != nil {
		return nil, err
	}

	userID, err := id.Unmarshal(regConfirmRequest.GetUserID())
	if err != nil {
		return nil, errors.Errorf("Unable to unmarshal user ID: %+v", err)
	}

	// Obtain signed client public key by passing to server
	signature, clientGwKey, err := s.handler.ConfirmRegistration(userID,
		regConfirmRequest.NonceSignedByClient.Signature, authState)

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the RegistrationConfirmation
	return &pb.RegistrationConfirmation{
		ClientSignedByServer: &messages.RSASignature{
			Signature: signature,
		},
		Error:            errMsg,
		ClientGatewayKey: clientGwKey,
	}, err
}

// PostPrecompResult sends final Message and AD precomputations.
func (s *Comms) PostPrecompResult(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	batchMsg := &pb.Batch{}
	err = ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, err
	}

	// Call the server handler to start a new round
	err = s.handler.PostPrecompResult(batchMsg.GetRound().GetID(),
		batchMsg.Slots, authState)
	return &messages.Ack{}, err
}

// FinishRealtime broadcasts to all nodes when the realtime is completed
func (s *Comms) FinishRealtime(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	roundInfoMsg := &pb.RoundInfo{}
	err = ptypes.UnmarshalAny(msg.Message, roundInfoMsg)
	if err != nil {
		return nil, err
	}

	err = s.handler.FinishRealtime(roundInfoMsg, authState)

	return &messages.Ack{}, err
}

// GetCompletedBatch should return a completed batch that the calling gateway
// hasn't gotten before
func (s *Comms) GetCompletedBatch(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*pb.Batch, error) {

	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	return s.handler.GetCompletedBatch(authState)
}

func (s *Comms) GetMeasure(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.RoundMetrics, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	roundInfoMsg := &pb.RoundInfo{}
	err = ptypes.UnmarshalAny(msg.Message, roundInfoMsg)
	if err != nil {
		return nil, err
	}

	rm, err := s.handler.GetMeasure(roundInfoMsg, authState)
	return rm, err
}

// Gateway -> Server unified polling
func (s *Comms) Poll(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.ServerPollResponse, error) {
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	//Unmarshall the any message to the message type needed
	pollMsg := &pb.ServerPoll{}
	err = ptypes.UnmarshalAny(msg.Message, pollMsg)
	if err != nil {
		return nil, err
	}

	return s.handler.Poll(pollMsg, authState)
}

func (s *Comms) RoundError(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable to handle reception of AuthenticatedMessage: %+v", err)
	}
	errMsg := &pb.RoundError{}
	err = ptypes.UnmarshalAny(msg.Message, errMsg)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, s.handler.RoundError(errMsg, authState)
}

func (s *Comms) SendRoundTripPing(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	//Marshall the any message to the message type needed
	roundTripPing := &pb.RoundTripPing{}
	err = ptypes.UnmarshalAny(msg.Message, roundTripPing)
	if err != nil {
		return nil, err
	}

	err = s.handler.SendRoundTripPing(roundTripPing, authState)
	return &messages.Ack{}, err
}
