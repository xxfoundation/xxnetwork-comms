////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains authorizer server gRPC endpoints

package authorizer

import (
	"context"

	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
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

// Authorizes a node to talk to permissioning
func (r *Comms) Authorize(ctx context.Context, auth *pb.AuthorizerAuth) (ack *messages.Ack, err error) {
	address, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return &messages.Ack{Error: err.Error()}, err
	}
	returned_err := r.handler.Authorize(auth, address)
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	return &messages.Ack{Error: errString}, returned_err
}
