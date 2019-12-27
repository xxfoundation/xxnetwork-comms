////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for registration functionality

package registration

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
func StartRegistrationServer(id, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {

	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Printf("Unable to start comms server: %+v", err)
	}

	registrationServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		pb.RegisterRegistrationServer(registrationServer.LocalServer, &registrationServer)
		pb.RegisterGenericServer(registrationServer.LocalServer, &registrationServer)

		// Register reflection service on gRPC server.
		reflection.Register(registrationServer.LocalServer)
		if err := registrationServer.LocalServer.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
		jww.INFO.Printf("Shutting down registration server listener:"+
			" %s", lis)
	}()

	return &registrationServer
}

type Handler interface {
	RegisterUser(registrationCode, pubKey string) (signature []byte, err error)
	GetCurrentClientVersion() (version string, err error)
	RegisterNode(ID []byte, ServerAddr, ServerTlsCert, GatewayAddr,
		GatewayTlsCert, RegistrationCode string) error
	PollNdf(ndfHash []byte, auth *connect.Auth) ([]byte, error)
}

type implementationFunctions struct {
	RegisterUser func(registrationCode, pubKey string) (signature []byte,
		err error)
	GetCurrentClientVersion func() (version string, err error)
	RegisterNode            func(ID []byte, ServerAddr, ServerTlsCert,
		GatewayAddr, GatewayTlsCert, RegistrationCode string) error
	PollNdf func(ndfHash []byte, auth *connect.Auth) ([]byte, error)
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

			RegisterUser: func(registrationCode,
				pubKey string) (signature []byte, err error) {
				warn(um)
				return nil, nil
			},
			GetCurrentClientVersion: func() (version string, err error) {
				warn(um)
				return "", nil
			},
			RegisterNode: func(ID []byte, ServerAddr, ServerTlsCert,
				GatewayAddr, GatewayTlsCert, RegistrationCode string) error {
				warn(um)
				return nil
			},
			PollNdf: func(ndfHash []byte, auth *connect.Auth) ([]byte, error) {
				warn(um)
				return nil, nil
			},
		},
	}
}

// Registers a user and returns a signed public key
func (s *Implementation) RegisterUser(registrationCode,
	pubKey string) (signature []byte, err error) {
	return s.Functions.RegisterUser(registrationCode, pubKey)
}

func (s *Implementation) GetCurrentClientVersion() (string, error) {
	return s.Functions.GetCurrentClientVersion()
}

func (s *Implementation) RegisterNode(ID []byte, ServerAddr, ServerTlsCert,
	GatewayAddr, GatewayTlsCert, RegistrationCode string) error {
	return s.Functions.RegisterNode(ID, ServerAddr, ServerTlsCert,
		GatewayAddr, GatewayTlsCert, RegistrationCode)
}

func (s *Implementation) PollNdf(ndfHash []byte, auth *connect.Auth) ([]byte, error) {
	return s.Functions.PollNdf(ndfHash, auth)
}
