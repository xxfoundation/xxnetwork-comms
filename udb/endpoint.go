////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains user discovery server gRPC endpoint wrappers
// When you add the udb server to mixmessages/mixmessages.proto and add the
// first function, a version of that goes here which calls the "handler"
// version of the function, with any mappings/wrappings necessary.

package udb

import (
	"context"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/messages"
)

// Handles validation of reverse-authentication tokens
func (u *Comms) AuthenticateToken(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	err := u.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &messages.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (u *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := u.GenerateToken()
	return &messages.AssignToken{
		Token: token,
	}, err
}

func (u *Comms) RegisterUser(ctx context.Context, msg *pb.UDBUserRegistration) (*messages.Ack, error) {
	return u.handler.RegisterUser(msg)
}

func (u *Comms) RemoveUser(ctx context.Context, msg *pb.FactRemovalRequest) (*messages.Ack, error) {
	return u.handler.RemoveUser(msg)
}

func (u *Comms) RegisterFact(ctx context.Context, msg *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error) {
	return u.handler.RegisterFact(msg)
}

func (u *Comms) ConfirmFact(ctx context.Context, msg *pb.FactConfirmRequest) (*messages.Ack, error) {
	return u.handler.ConfirmFact(msg)
}

func (u *Comms) RemoveFact(ctx context.Context, msg *pb.FactRemovalRequest) (*messages.Ack, error) {
	return u.handler.RemoveFact(msg)
}
