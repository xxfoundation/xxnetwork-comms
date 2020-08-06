///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package interconnect

import (
	"errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/reflection"
	"net"
	"runtime/debug"
)

// Starts a new server on the localHost:port specified by port
// and a callback interface for interconnect operations
// with given path to public and private key for TLS connection
func StartCMixInterconnect(id *id.ID, port string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {

	addr := net.JoinHostPort("0.0.0.0", port)

	pc, lis, err := connect.StartCommServer(id, addr,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	CMixInterconnect := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		// Register GRPC services to the listening address
		RegisterInterconnectServer(CMixInterconnect.LocalServer, &CMixInterconnect)
		//messages.RegisterGenericServer(CMixInterconnect.LocalServer, &CMixInterconnect)

		// Register reflection service on gRPC server.
		reflection.Register(CMixInterconnect.LocalServer)
		if err := CMixInterconnect.LocalServer.Serve(lis); err != nil {
			jww.FATAL.Panicf("Failed to serve: %+v",
				errors.New(err.Error()))
		}
		jww.INFO.Printf("Shutting down node server listener: %s", lis)
	}()

	return &CMixInterconnect

}

// Server object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
}

type Handler interface {
	// Interconnect interface for getting the NDF
	GetNDF() (*NDF, error)
}

type implementationFunctions struct {
	GetNDF func() (*NDF, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// Below is the Implementation implementation, which calls the
// function matching the variable in the structure.

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
			GetNDF: func() (*NDF, error) {
				warn(um)
				return &NDF{
					Ndf: []byte("hello world"),
				}, nil
			},
		},
	}
}

// Interconnect Interface for getting an NDF
func (s *Implementation) GetNDF() (*NDF, error) {
	return s.Functions.GetNDF()
}
