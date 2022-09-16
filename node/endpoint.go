////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains server gRPC endpoints

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// Handle a Broadcasted Ask Online event
func (s *Comms) AskOnline(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	auth, err := s.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//return an error if the connection is not authenticated
	if !auth.IsAuthenticated {
		return &messages.Ack{}, connect.AuthError(auth.Sender.GetId())
	}

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
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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

// Handle a Phase event
func (s *Comms) PostPhase(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack,
	error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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

	authState, err := s.AuthenticatedReceiver(authMsg, server.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Verify the message authentication
	return s.handler.StreamPostPhase(server, authState)
}

// GetBufferInfo returns buffer size (number of completed precomputations)
func (s *Comms) GetRoundBufferInfo(ctx context.Context,
	msg *messages.AuthenticatedMessage) (
	*pb.RoundBufferInfo, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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
func (s *Comms) RequestClientKey(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*pb.SignedKeyResponse, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Marshall the any message to the message type needed
	nonceRequest := &pb.SignedClientKeyRequest{}
	err = ptypes.UnmarshalAny(msg.Message, nonceRequest)
	if err != nil {
		return nil, err
	}

	// Obtain the nonce by passing to server
	nonce, err := s.handler.RequestClientKey(nonceRequest, authState)
	if err != nil {

	}

	// Return the NonceMessage
	return nonce, err
}

// PostPrecompResult sends final Message and AD precomputations.
func (s *Comms) PostPrecompResult(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {

	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	batchMsg := &pb.PostPrecompResult{}
	err = ptypes.UnmarshalAny(msg.Message, batchMsg)
	if err != nil {
		return nil, err
	}

	// Call the server handler to start a new round
	err = s.handler.PostPrecompResult(batchMsg.RoundId,
		batchMsg.NumSlots, authState)
	return &messages.Ack{}, err
}

func (s *Comms) GetMeasure(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.RoundMetrics, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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
	authState, err := s.AuthenticatedReceiver(msg, ctx)
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

// Server -> Gateway permissioning address
func (s *Comms) GetPermissioningAddress(context.Context, *messages.Ping) (*pb.StrAddress, error) {
	ip, err := s.handler.GetPermissioningAddress()
	if err != nil {
		return nil, err
	}
	return &pb.StrAddress{Address: ip}, nil
}

// Server -> Server initiating multi-party round DH key generation
func (s *Comms) StartSharePhase(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}
	//Marshall the any message to the message type needed
	startShare := &pb.RoundInfo{}
	err = ptypes.UnmarshalAny(msg.Message, startShare)
	if err != nil {
		return nil, err
	}

	err = s.handler.StartSharePhase(startShare, authState)
	return &messages.Ack{}, err

}

// Server -> Server passing state of multi-party round DH key generation
func (s *Comms) SharePhaseRound(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Marshall the any message to the message type needed
	sharePiece := &pb.SharePiece{}
	err = ptypes.UnmarshalAny(msg.Message, sharePiece)
	if err != nil {
		return nil, err
	}

	err = s.handler.SharePhaseRound(sharePiece, authState)
	return &messages.Ack{}, err
}

// Server -> Server sending multi-party round DH final key
func (s *Comms) ShareFinalKey(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Verify the message authentication
	authState, err := s.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Marshall the any message to the message type needed
	sharePiece := &pb.SharePiece{}
	err = ptypes.UnmarshalAny(msg.Message, sharePiece)
	if err != nil {
		return nil, err
	}

	err = s.handler.ShareFinalKey(sharePiece, authState)
	return &messages.Ack{}, err
}
