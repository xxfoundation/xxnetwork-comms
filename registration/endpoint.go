////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server gRPC endpoints

package registration

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

// RegisterUser event handler which registers a user with the platform
func (s *RegistrationComms) RegisterUser(ctx context.Context, msg *pb.UserRegistration) (
	*pb.UserRegistrationConfirmation, error) {
	// Obtain the signed key by passing to registration server
	pubKey := msg.GetClient()
	hash, R, S, err := s.handler.RegisterUser(msg.
		GetRegistrationCode(), pubKey.GetY(), pubKey.GetP(),
		pubKey.GetQ(), pubKey.GetG())
	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
		err = errors.New(err.Error())
	}

	// Return the confirmation message
	return &pb.UserRegistrationConfirmation{
		ClientSignedByServer: &pb.DSASignature{
			Hash: hash,
			R:    R,
			S:    S,
		},
		Error: errMsg,
	}, err
}

// Handle a node registration event
func (s *RegistrationComms) RegisterNode(ctx context.Context, msg *pb.NodeRegistration) (
	*pb.Ack, error) {
	info, _ := peer.FromContext(ctx)

	err := s.handler.RegisterNode(msg.GetID(), msg.GetNodeCSR(), msg.GetGatewayTLSCert(), msg.GetRegistrationCode(), info.Addr.String())
	return &pb.Ack{}, err
}
