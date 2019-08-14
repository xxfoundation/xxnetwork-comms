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
	RegisterNode(ID []byte,
		ServerTlsCert, GatewayTlsCert, RegistrationCode, Addr string) error
}

type implementationFunctions struct {
	RegisterUser func(registrationCode, pubKey string) (signature []byte,
		err error)
	RegisterNode func(ID []byte,
		ServerTlsCert, GatewayTlsCert, RegistrationCode, Addr string) error
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
func NewImplementation() Handler {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return Handler(&Implementation{
		Functions: implementationFunctions{

			RegisterUser: func(registrationCode,
				pubKey string) (signature []byte, err error) {
				warn(um)
				return nil, nil
			},
			RegisterNode: func(ID []byte,
				ServerTlsCert, GatewayTlsCert, RegistrationCode, Addr string) error {
				warn(um)
				return nil
			},
		},
	})
}

// Registers a user and returns a signed public key
func (s *Implementation) RegisterUser(registrationCode,
	pubKey string) (signature []byte, err error) {
	return s.Functions.RegisterUser(registrationCode, pubKey)
}

func (s *Implementation) RegisterNode(ID []byte,
	ServerTlsCert, GatewayTlsCert, RegistrationCode, Addr string) error {
	return s.Functions.RegisterNode(ID, ServerTlsCert, GatewayTlsCert, RegistrationCode, Addr)

}
