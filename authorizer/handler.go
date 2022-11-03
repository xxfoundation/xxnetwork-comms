////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for authorizer functionality

package authorizer

import (
	"runtime/debug"

	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
)

// Authorizer object used to implement
// endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
	*pb.UnimplementedAuthorizerServer
	*messages.UnimplementedGenericServer
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartAuthorizerServer(id *id.ID, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {

	pc, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock, nil)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	authorizerServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}
	pb.RegisterAuthorizerServer(authorizerServer.GetServer(), &authorizerServer)
	messages.RegisterGenericServer(authorizerServer.GetServer(), &authorizerServer)

	pc.Serve()
	return &authorizerServer
}

type Handler interface {
	Authorize(auth *pb.AuthorizerAuth, ipAddr string) (err error)
	RequestCert(msg *pb.AuthorizerCertRequest) (*messages.Ack, error)
	SendEABCredentialRequest(msg *pb.EABCredentialRequest) (*pb.EABCredentialResponse, error)
}

type implementationFunctions struct {
	Authorize                func(auth *pb.AuthorizerAuth, ipAddr string) (err error)
	RequestCert              func(msg *pb.AuthorizerCertRequest) (*messages.Ack, error)
	SendEABCredentialRequest func(msg *pb.EABCredentialRequest) (*pb.EABCredentialResponse, error)
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

			Authorize: func(auth *pb.AuthorizerAuth, ipAddr string) (err error) {
				warn(um)
				return nil
			},
			RequestCert: func(msg *pb.AuthorizerCertRequest) (*messages.Ack, error) {
				warn(um)
				return &messages.Ack{}, nil
			},
			SendEABCredentialRequest: func(msg *pb.EABCredentialRequest) (*pb.EABCredentialResponse, error) {
				warn(um)
				return &pb.EABCredentialResponse{}, nil
			},
		},
	}
}

// Authorizes a node to talk to permissioning
func (s *Implementation) Authorize(auth *pb.AuthorizerAuth, ipAddr string) (err error) {
	return s.Functions.Authorize(auth, ipAddr)
}

// Request a signed certificate for HTTPS
func (s *Implementation) RequestCert(msg *pb.AuthorizerCertRequest) (*messages.Ack, error) {
	return s.Functions.RequestCert(msg)
}

// Request ACME key for HTTPS
func (s *Implementation) SendEABCredentialRequest(msg *pb.EABCredentialRequest) (*pb.EABCredentialResponse, error) {
	return s.Functions.SendEABCredentialRequest(msg)
}
