////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package notificationBot

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Handler interface for the Notification Bot
type Handler interface {
	// RegisterForNotifications event handler which registers a client with the notification bot
	RegisterForNotifications(clientToken []byte, auth *connect.Auth) error
	// UnregisterForNotifications event handler which unregisters a client with the notification bot
	UnregisterForNotifications(auth *connect.Auth) error
	// Ping gateway to ask for users to notify
	PollForNotifications(auth *connect.Auth) ([]string, error)
}

// NotificationBot object used to implement
// endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartNotificationBot(id, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {

	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	notificationBot := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		pb.RegisterNotificationBotServer(notificationBot.LocalServer, &notificationBot)
		pb.RegisterGenericServer(notificationBot.LocalServer, &notificationBot)

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
	RegisterForNotifications   func(clientToken []byte, auth *connect.Auth) error
	UnregisterForNotifications func(auth *connect.Auth) error
	PollForNotifications       func(auth *connect.Auth) ([]string, error)
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

			RegisterForNotifications: func(clientToken []byte, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			UnregisterForNotifications: func(auth *connect.Auth) error {
				warn(um)
				return nil
			},
			PollForNotifications: func(auth *connect.Auth) ([]string, error) {
				warn(um)
				return nil, nil
			},
		},
	}
}

// RegisterForNotifications event handler which registers a client with the notification bot
func (s *Implementation) RegisterForNotifications(clientToken []byte, auth *connect.Auth) error {
	return s.Functions.RegisterForNotifications(clientToken, auth)
}

// UnregisterForNotifications event handler which unregisters a client with the notification bot
func (s *Implementation) UnregisterForNotifications(auth *connect.Auth) error {
	return s.Functions.UnregisterForNotifications(auth)
}

// Ping gateway to ask for users to notify
func (s *Implementation) PollForNotifications(auth *connect.Auth) ([]string, error) {
	return s.Functions.PollForNotifications(auth)
}
