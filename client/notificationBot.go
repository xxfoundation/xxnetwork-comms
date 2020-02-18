////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> notificationBot functionality

package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// Client -> NotificationBot
func (c *Comms) RegisterForNotifications(host *connect.Host,
	message *pb.NotificationToken) (*pb.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		authMsg, err := c.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewNotificationBotClient(conn).RegisterForNotifications(ctx,
			authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending RegisterForNotification message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)

}

// Client -> NotificationBot
func (c *Comms) UnregisterForNotifications(host *connect.Host) (*pb.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		authMsg, err := c.PackAuthenticatedMessage(&pb.Ping{}, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewNotificationBotClient(conn).UnregisterForNotifications(ctx,
			authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending UnregisterForNotification message")
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)

}
