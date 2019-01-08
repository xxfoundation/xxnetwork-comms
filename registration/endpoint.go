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

// Handle a Register User event
func (s *server) RegisterUser(ctx context.Context, msg *pb.RegisterUserMessage) (
	*pb.ConfirmRegisterUserMessage, error) {

	signedKey, err := registrationHandler.RegisterUser(msg.
		RegistrationCode, msg.Email, msg.Password, msg.PublicKey)

	return &pb.ConfirmRegisterUserMessage{
		SignedPublicKey: signedKey,
		Error:           err.Error(),
	}, err
}
