////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains registration server gRPC endpoints

package clientregistrar

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (r *Comms) AuthenticateToken(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	err := r.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &messages.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (r *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := r.GenerateToken()
	return &messages.AssignToken{
		Token: token,
	}, err
}

// RegisterUser event handler which registers a user with the platform
func (r *Comms) RegisterUser(ctx context.Context, msg *pb.ClientRegistration) (
	*pb.SignedClientRegistrationConfirmations, error) {
	// Obtain the signed key by passing to registration server
	confirmationMessage, err := r.handler.RegisterUser(msg)
	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
		err = errors.New(err.Error())
	}
	confirmationMessage.Error = errMsg
	// Return the confirmation message
	return confirmationMessage, err
}
