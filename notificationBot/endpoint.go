////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains notificationBot gRPC endpoints

package notificationBot

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (nb *Comms) AuthenticateToken(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	err := nb.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &pb.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (nb *Comms) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	token, err := nb.GenerateToken()
	return &pb.AssignToken{
		Token: token,
	}, err
}

// RegisterForNotifications event handler which registers a client with the notification bot
func (nb *Comms) RegisterForNotifications(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	//Check the authState of the message
	authState := nb.AuthenticatedReceiver(msg)

	notificationToken := &pb.NotificationToken{}
	err := ptypes.UnmarshalAny(msg.Message, notificationToken)
	if err != nil {
		return nil, err
	}

	err = nb.handler.RegisterForNotifications(notificationToken.Token, authState)
	// Obtain the error message, if any
	if err != nil {
		err = errors.New(err.Error())
	}

	// Return the confirmation message
	return &pb.Ack{}, err
}

// UnregisterForNotifications event handler which unregisters a client with the notification bot
func (nb *Comms) UnregisterForNotifications(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	// Check the authState of the message
	authState := nb.AuthenticatedReceiver(msg)

	err := nb.handler.UnregisterForNotifications(authState)
	// Obtain the error message, if any
	if err != nil {
		err = errors.New(err.Error())
	}

	// Return the confirmation message
	return &pb.Ack{}, err
}

// Ping gateway to ask for users to notify
func (nb *Comms) PollForNotifications(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.IDList, error) {
	// Check the authState of the message
	authState := nb.AuthenticatedReceiver(msg)

	ids, err := nb.handler.PollForNotifications(authState)
	return &pb.IDList{
		IDs: ids,
	}, err
}
