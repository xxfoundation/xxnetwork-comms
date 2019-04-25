////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server GRPC endpoints

package registration

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// RegisterUser event handler which registers a user with the platform
func (s *server) RegisterUser(ctx context.Context, msg *pb.UserRegistration) (
	*pb.UserRegistrationConfirmation, error) {

	// Obtain the signed key by passing to registration server
	pubKey := msg.GetClient()
	hash, R, S, err := registrationHandler.RegisterUser(msg.
		GetRegistrationCode(), pubKey.GetY(), pubKey.GetP(),
		pubKey.GetQ(), pubKey.GetG())

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
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
