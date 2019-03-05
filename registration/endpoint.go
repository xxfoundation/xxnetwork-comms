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

// Handle a RegisterUser event
func (s *server) RegisterUser(ctx context.Context, msg *pb.RegisterUserMessage) (
	*pb.ConfirmRegisterUserMessage, error) {

	// Obtain the signed key by passing to registration server
	hash, R, S, err := registrationHandler.RegisterUser(msg.
		GetRegistrationCode(), msg.GetY(), msg.GetP(), msg.GetQ(), msg.GetG())

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the confirmation message
	return &pb.ConfirmRegisterUserMessage{
		Hash:  hash,
		R:     R,
		S:     S,
		Error: errMsg,
	}, err
}
