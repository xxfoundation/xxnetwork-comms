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

type RegistrationHandler interface {
	RegisterUser(registrationCode, email, password string,
		publicKey []byte) ([]byte, error)
}

type implementationFunctions struct {
	RegisterUser func(registrationCode, email, password string,
		publicKey []byte) ([]byte, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
func NewImplementation() RegistrationHandler {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%v", debug.Stack())
	}
	return RegistrationHandler(&Implementation{
		Functions: implementationFunctions{
			RegisterUser: func(registrationCode, email, password string,
				publicKey []byte) ([]byte, error) {
				warn(um)
				return nil, nil
			},
		},
	})
}

// Registers a user and returns a signed public key
func (s *Implementation) RegisterUser(registrationCode, email, password string,
	publicKey []byte) ([]byte, error) {
	return s.Functions.RegisterUser(registrationCode, email, password, publicKey)
}
