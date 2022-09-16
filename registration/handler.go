////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for registration functionality

package registration

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

// Registration object used to implement
// endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
	*pb.UnimplementedRegistrationServer
	*messages.UnimplementedGenericServer
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartRegistrationServer(id *id.ID, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte, preloadedHosts []*connect.Host) *Comms {

	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock, preloadedHosts)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	registrationServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		pb.RegisterRegistrationServer(registrationServer.LocalServer, &registrationServer)
		messages.RegisterGenericServer(registrationServer.LocalServer, &registrationServer)

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
	RegisterUser(msg *pb.ClientRegistration) (confirmation *pb.SignedClientRegistrationConfirmations, err error)
	RegisterNode(salt []byte, serverAddr, serverTlsCert, gatewayAddr,
		gatewayTlsCert, registrationCode string) error
	PollNdf(ndfHash []byte) (*pb.NDF, error)
	Poll(msg *pb.PermissioningPoll, auth *connect.Auth) (*pb.
		PermissionPollResponse, error)
	CheckRegistration(msg *pb.RegisteredNodeCheck) (*pb.RegisteredNodeConfirmation, error)
}

type implementationFunctions struct {
	RegisterUser func(msg *pb.ClientRegistration) (confirmation *pb.SignedClientRegistrationConfirmations, err error)
	RegisterNode func(salt []byte, serverAddr, serverTlsCert, gatewayAddr,
		gatewayTlsCert, registrationCode string) error
	PollNdf           func(ndfHash []byte) (*pb.NDF, error)
	Poll              func(msg *pb.PermissioningPoll, auth *connect.Auth) (*pb.PermissionPollResponse, error)
	CheckRegistration func(msg *pb.RegisteredNodeCheck) (*pb.RegisteredNodeConfirmation, error)
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

			RegisterUser: func(msg *pb.ClientRegistration) (confirmation *pb.SignedClientRegistrationConfirmations, err error) {
				warn(um)
				return &pb.SignedClientRegistrationConfirmations{}, nil
			},
			RegisterNode: func(salt []byte, serverAddr, serverTlsCert, gatewayAddr,
				gatewayTlsCert, registrationCode string) error {
				warn(um)
				return nil
			},
			PollNdf: func(ndfHash []byte) (*pb.NDF, error) {
				warn(um)
				return nil, nil
			},
			Poll: func(msg *pb.PermissioningPoll, auth *connect.Auth) (*pb.PermissionPollResponse, error) {
				warn(um)
				return &pb.PermissionPollResponse{}, nil
			},
			CheckRegistration: func(msg *pb.RegisteredNodeCheck) (*pb.
				RegisteredNodeConfirmation, error) {

				warn(um)
				return &pb.RegisteredNodeConfirmation{}, nil
			},
		},
	}
}

// Registers a user and returns a signed public key
func (s *Implementation) RegisterUser(msg *pb.ClientRegistration) (confirmation *pb.SignedClientRegistrationConfirmations, err error) {
	return s.Functions.RegisterUser(msg)
}

func (s *Implementation) RegisterNode(salt []byte, serverAddr, serverTlsCert, gatewayAddr,
	gatewayTlsCert, registrationCode string) error {
	return s.Functions.RegisterNode(salt, serverAddr, serverTlsCert, gatewayAddr,
		gatewayTlsCert, registrationCode)
}

func (s *Implementation) PollNdf(ndfHash []byte) (*pb.NDF, error) {
	return s.Functions.PollNdf(ndfHash)
}

func (s *Implementation) Poll(msg *pb.PermissioningPoll, auth *connect.Auth) (*pb.PermissionPollResponse, error) {
	return s.Functions.Poll(msg, auth)
}

func (s *Implementation) CheckRegistration(msg *pb.RegisteredNodeCheck) (*pb.
	RegisteredNodeConfirmation, error) {
	return s.Functions.CheckRegistration(msg)
}
