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

// Registration object used to implement
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
		//Change func calls
		pb.RegisterNotificationBotServer(notificationBot.LocalServer, &notificationBot)
		//Need this?
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

type Handler interface {
	RegisterForNotifications(clientToken []byte, auth *connect.Auth) error
	UnregisterForNotifications(auth *connect.Auth) error
}

type implementationFunctions struct {
	RegisterForNotifications   func(clientToken []byte, auth *connect.Auth) error
	UnregisterForNotifications func(auth *connect.Auth) error
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
		},
	}
}
