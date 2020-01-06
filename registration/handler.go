////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for registration functionality

package registration

import (
	jww "github.com/spf13/jwalterweatherman"
	"runtime/debug"
)

type Handler interface {
	RegisterUser(registrationCode, pubKey string) (signature []byte, err error)
	GetCurrentClientVersion() (version string, err error)
	RegisterNode(ID []byte, ServerAddr, ServerTlsCert, GatewayAddr,
		GatewayTlsCert, RegistrationCode string) error
	PollNdf(ndfHash []byte) ([]byte, error)
}

type implementationFunctions struct {
	RegisterUser func(registrationCode, pubKey string) (signature []byte,
		err error)
	GetCurrentClientVersion func() (version string, err error)
	RegisterNode            func(ID []byte, ServerAddr, ServerTlsCert,
		GatewayAddr, GatewayTlsCert, RegistrationCode string) error
	PollNdf func(ndfHash []byte) ([]byte, error)
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
			PollNdf: func(ndfHash []byte) ([]byte, error) {
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

func (s *Implementation) PollNdf(ndfHash []byte) ([]byte, error) {
	return s.Functions.PollNdf(ndfHash)
}
