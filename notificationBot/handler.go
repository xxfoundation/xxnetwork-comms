////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for notificationBot functionality

package notificationBot

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Handler interface for the Notification Bot
type Handler interface {
	// RegisterForNotifications event handler which registers a client with the notification bot
	RegisterForNotifications(msg *pb.NotificationRegisterRequest) error
	// UnregisterForNotifications event handler which unregisters a client with the notification bot
	UnregisterForNotifications(msg *pb.NotificationUnregisterRequest) error
	// ReceiveNotificationBatch receives the batch of notification data from gateway.
	ReceiveNotificationBatch(notifBatch *pb.NotificationBatch, auth *connect.Auth) error
}

// NotificationBot object used to implement
// endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
	*pb.UnimplementedNotificationBotServer
	*messages.UnimplementedGenericServer
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartNotificationBot(id *id.ID, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {

	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock, nil)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	notificationBot := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		pb.RegisterNotificationBotServer(notificationBot.LocalServer, &notificationBot)
		messages.RegisterGenericServer(notificationBot.LocalServer, &notificationBot)

		// Register reflection service on gRPC server.
		reflection.Register(notificationBot.LocalServer)
		if err := notificationBot.LocalServer.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
		jww.INFO.Printf("Shutting down registration server listener:"+
			" %s", lis)
	}()

	return &notificationBot
}

// Handler implementation for the NotificationBot
type implementationFunctions struct {
	RegisterForNotifications   func(request *pb.NotificationRegisterRequest) error
	UnregisterForNotifications func(request *pb.NotificationUnregisterRequest) error
	ReceiveNotificationBatch   func(notifBatch *pb.NotificationBatch, auth *connect.Auth) error
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{

			RegisterForNotifications: func(request *pb.NotificationRegisterRequest) error {
				warn(um)
				return nil
			},
			UnregisterForNotifications: func(request *pb.NotificationUnregisterRequest) error {
				warn(um)
				return nil
			},
			ReceiveNotificationBatch: func(notifBatch *pb.NotificationBatch, auth *connect.Auth) error {
				warn(um)
				return nil
			},
		},
	}
}

// RegisterForNotifications event handler which registers a client with the notification bot
func (s *Implementation) RegisterForNotifications(request *pb.NotificationRegisterRequest) error {
	return s.Functions.RegisterForNotifications(request)
}

// UnregisterForNotifications event handler which unregisters a client with the notification bot
func (s *Implementation) UnregisterForNotifications(request *pb.NotificationUnregisterRequest) error {
	return s.Functions.UnregisterForNotifications(request)
}

// ReceiveNotificationBatch receives the batch of notification data from gateway.
func (s *Implementation) ReceiveNotificationBatch(notifBatch *pb.NotificationBatch,
	auth *connect.Auth) error {
	return s.Functions.ReceiveNotificationBatch(notifBatch, auth)
}
