///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains notificationBot gRPC endpoints

package notificationBot

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (nb *Comms) AuthenticateToken(_ context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	err := nb.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &messages.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (nb *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := nb.GenerateToken()
	return &messages.AssignToken{
		Token: token,
	}, err
}

// RegisterForNotifications event handler which registers a client with the notification bot
func (nb *Comms) RegisterForNotifications(_ context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Check the authState of the message
	authState, err := nb.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Failed to handle reception of AuthenticatedMessage: %+v", err)
	}

	req := &pb.NotificationRegisterRequest{}
	err = ptypes.UnmarshalAny(msg.Message, req)
	if err != nil {
		return nil, err
	}

	err = nb.handler.RegisterForNotifications(req, authState)

	// Return the confirmation message
	return &messages.Ack{}, err
}

// UnregisterForNotifications event handler which unregisters a client with the notification bot
func (nb *Comms) UnregisterForNotifications(_ context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Check the authState of the message
	authState, err := nb.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Failed to handle reception of AuthenticatedMessage: %+v", err)
	}

	req := &pb.NotificationUnregisterRequest{}
	err = ptypes.UnmarshalAny(msg.Message, req)
	if err != nil {
		return nil, err
	}

	err = nb.handler.UnregisterForNotifications(req, authState)

	// Return the confirmation message
	return &messages.Ack{}, err
}

func (nb *Comms) ReceiveNotificationBatch(_ context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Check the authState of the message
	authState, err := nb.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Failed to handle reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshal notification data
	notificationBatch := &pb.NotificationBatch{}
	err = ptypes.UnmarshalAny(msg.Message, notificationBatch)
	if err != nil {
		return nil, err
	}

	err = nb.handler.ReceiveNotificationBatch(notificationBatch, authState)

	return &messages.Ack{}, err
}
