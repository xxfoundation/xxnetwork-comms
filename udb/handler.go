///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains callback interface for registration functionality

package udb

import (
	//	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	//	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	//	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Registration object used to implement
// endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler // an object that implements the interface below, which
	// has all the functions called by endpoint.go
}

// StartServer starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartServer(id *id.ID, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	/*
		pc, lis, err := connect.StartCommServer(id, localServer,
			certPEMblock, keyPEMblock)
		if err != nil {
			jww.FATAL.Panicf("Unable to start comms server: %+v", err)
		}

		udbServer := Comms{
			ProtoComms: pc,
			handler:    handler,
		}

		go func() {
			pb.RegisterUDBServer(udbServer.LocalServer,
				&udbServer)
			messages.RegisterGenericServer(udbServer.LocalServer,
				&udbServer)

			// Register reflection service on gRPC server.
			reflection.Register(udbServer.LocalServer)
			if err := udbServer.LocalServer.Serve(lis); err != nil {
				err = errors.New(err.Error())
				jww.FATAL.Panicf("Failed to serve: %+v", err)
			}
			jww.INFO.Printf("Shutting down registration server listener:"+
				" %s", lis)
		}()
		return &udbServer
	*/
	return nil
}

// Handler is the interface udb has to implement to integrate with the comms
// library properly.
type Handler interface {
	// ClientCall inside UDB needs to implement this interface.
	ClientCall(msg *pb.PermissioningPoll, auth *connect.Auth,
		serverAddress string) (*pb.PermissionPollResponse, error)
	RegisterUser(registration *pb.UDBUserRegistration) pb.UserRegistrationResponse
	RegisterFact(request *pb.FactRegisterRequest) pb.FactRegisterResponse
	ConfirmFact(request *pb.FactConfirmRequest) pb.FactConfirmResponse
	RemoveFact(request *pb.FactRemovalRequest) pb.FactRemovalResponse
}

// implementationFunctions are the actual implementations of
type implementationFunctions struct {
	// This is the function "implementation" -- inside UDB we will
	// set this to be the UDB version of the function. By default
	// it's a dummy function that returns nothing (see NewImplementation
	// below).
	ClientCall func(msg *pb.PermissioningPoll, auth *connect.Auth,
		serverAddress string) (*pb.PermissionPollResponse, error)
	RegisterUser func(registration *pb.UDBUserRegistration) pb.UserRegistrationResponse
	RegisterFact func(request *pb.FactRegisterRequest) pb.FactRegisterResponse
	ConfirmFact  func(request *pb.FactConfirmRequest) pb.FactConfirmResponse
	RemoveFact   func(request *pb.FactRemovalRequest) pb.FactRemovalResponse
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
// Inside UDB, you would call this, then set "ClientCall" to your
// own UDB version of the function.
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			ClientCall: func(msg *pb.PermissioningPoll,
				auth *connect.Auth,
				serverAddress string) (
				*pb.PermissionPollResponse, error) {
				warn(um)
				return &pb.PermissionPollResponse{}, nil
			},
			RegisterUser: func(registration *pb.UDBUserRegistration) pb.UserRegistrationResponse {
				warn(um)
				return pb.UserRegistrationResponse{}
			},
			RegisterFact: func(request *pb.FactRegisterRequest) pb.FactRegisterResponse {
				warn(um)
				return pb.FactRegisterResponse{}
			},
			ConfirmFact: func(request *pb.FactConfirmRequest) pb.FactConfirmResponse {
				warn(um)
				return pb.FactConfirmResponse{}
			},
			RemoveFact: func(request *pb.FactRemovalRequest) pb.FactRemovalResponse {
				warn(um)
				return pb.FactRemovalResponse{}
			},
		},
	}
}

// ClientCall is called by the ClientCall in endpoint.go, which then calls
// the function inside the implementationFunctions struct. It's made to
// implement the interface.
func (s *Implementation) ClientCall(msg *pb.PermissioningPoll,
	auth *connect.Auth, serverAddress string) (
	*pb.PermissionPollResponse, error) {
	return s.Functions.ClientCall(msg, auth, serverAddress)
}

func (s *Implementation) RegisterUser(registration *pb.UDBUserRegistration) pb.UserRegistrationResponse {
	return s.Functions.RegisterUser(registration)
}

func (s *Implementation) RegisterFact(request *pb.FactRegisterRequest) pb.FactRegisterResponse {
	return s.Functions.RegisterFact(request)
}

func (s *Implementation) ConfirmFact(request *pb.FactConfirmRequest) pb.FactConfirmResponse {
	return s.Functions.ConfirmFact(request)
}

func (s *Implementation) RemoveFact(request *pb.FactRemovalRequest) pb.FactRemovalResponse {
	return s.Functions.RemoveFact(request)
}
