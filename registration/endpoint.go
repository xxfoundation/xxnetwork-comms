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
	id := msg.GetNodeId()
	nodeTLSCert := msg.GetNodeTLSCert()
	gatewayTLSCert := msg.GetGatewayTLSCert()
	regCode := msg.GetRegistrationCode()
	s.handler.RegisterNode(id, nodeTLSCert, gatewayTLSCert, regCode)
	return &pb.Ack{}, nil
}
