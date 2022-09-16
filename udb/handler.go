////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for registration functionality

package udb

import (
	"github.com/pkg/errors"
	//	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc/reflection"

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
	*pb.UnimplementedUDBServer
	*messages.UnimplementedGenericServer
}

// StartServer starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartServer(id *id.ID, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock, nil)
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
	return nil
}

// Handler is the interface udb has to implement to integrate with the comms
// library properly.
type Handler interface {
	// RegisterUser handles registering a user into the database
	RegisterUser(registration *pb.UDBUserRegistration) (*messages.Ack, error)
	// RemoveUser deletes this user registration and blocks anyone from ever
	// registering under that username again.
	// The fact removal request must be for the username or it will not work.
	RemoveUser(request *pb.FactRemovalRequest) (*messages.Ack, error)
	// RegisterFact handles registering a fact into the database
	RegisterFact(msg *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error)
	// ConfirmFact checks a Fact against the Fact database
	ConfirmFact(msg *pb.FactConfirmRequest) (*messages.Ack, error)
	// RemoveFact deletes a fact from its associated ID.
	// You cannot RemoveFact on a username. Callers must RemoveUser and reregister.
	RemoveFact(request *pb.FactRemovalRequest) (*messages.Ack, error)
}

// implementationFunctions are the actual implementations of
type implementationFunctions struct {
	// This is the function "implementation" -- inside UDB we will
	// set this to be the UDB version of the function. By default
	// it's a dummy function that returns nothing (see NewImplementation
	// below).

	// RegisterUser handles registering a user into the database
	RegisterUser func(registration *pb.UDBUserRegistration) (*messages.Ack, error)
	// RemoveUser deletes this user registration and blocks anyone from ever
	// registering under that username again.
	// The fact removal request must be for the username or it will not work.
	RemoveUser func(request *pb.FactRemovalRequest) (*messages.Ack, error)
	// RegisterFact handles registering a fact into the database
	RegisterFact func(request *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error)
	// ConfirmFact checks a Fact against the Fact database
	ConfirmFact func(request *pb.FactConfirmRequest) (*messages.Ack, error)
	// RemoveFact deletes a fact from its associated ID.
	// You cannot RemoveFact on a username. Callers must RemoveUser and reregister.
	RemoveFact func(request *pb.FactRemovalRequest) (*messages.Ack, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
// Inside UDB, you would call this, then set all functions to your
// own UDB version of the function.
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			// Stub for RegisterUser which returns a blank message and prints a warning
			RegisterUser: func(registration *pb.UDBUserRegistration) (*messages.Ack, error) {
				warn(um)
				return &messages.Ack{}, nil
			},
			// Stub for RemoveUser which returns a blank message and prints a warning
			RemoveUser: func(request *pb.FactRemovalRequest) (*messages.Ack, error) {
				warn(um)
				return &messages.Ack{}, nil
			},
			// Stub for RegisterFact which returns a blank message and prints a warning
			RegisterFact: func(request *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error) {
				warn(um)
				return &pb.FactRegisterResponse{}, nil
			},
			// Stub for ConfirmFact which returns a blank message and prints a warning
			ConfirmFact: func(request *pb.FactConfirmRequest) (*messages.Ack, error) {
				warn(um)
				return &messages.Ack{}, nil
			},
			// Stub for RemoveFact which returns a blank message and prints a warning
			RemoveFact: func(request *pb.FactRemovalRequest) (*messages.Ack, error) {
				warn(um)
				return &messages.Ack{}, nil
			},
		},
	}
}

// RegisterUser is called by the RegisterUser in endpoint.go. It calls the corresponding function in the interface.
func (s *Implementation) RegisterUser(registration *pb.UDBUserRegistration) (*messages.Ack, error) {
	return s.Functions.RegisterUser(registration)
}

// RemoveUser is called by the RemoveUser in endpoint.go. It calls the corresponding function in the interface.
func (s *Implementation) RemoveUser(request *pb.FactRemovalRequest) (*messages.Ack, error) {
	return s.Functions.RemoveUser(request)
}

// RegisterFact is called by the RegisterFact in endpoint.go. It calls the corresponding function in the interface.
func (s *Implementation) RegisterFact(request *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error) {
	return s.Functions.RegisterFact(request)
}

// ConfirmFact is called by the ConfirmFact in endpoint.go. It calls the corresponding function in the interface.
func (s *Implementation) ConfirmFact(request *pb.FactConfirmRequest) (*messages.Ack, error) {
	return s.Functions.ConfirmFact(request)
}

// RemoveFact is called by the RemoveFact in endpoint.go. It calls the corresponding function in the interface.
func (s *Implementation) RemoveFact(request *pb.FactRemovalRequest) (*messages.Ack, error) {
	return s.Functions.RemoveFact(request)
}
