////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Server implementation, interface & starter function

package server

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"runtime/debug"
)

// Comms object bundles low-level connect.ProtoComms,
// and the endpoint Handler interface.
type Comms struct {
	*connect.ProtoComms
	handler Handler
	*pb.UnimplementedRemoteSyncServer
	*messages.UnimplementedGenericServer
}

// Handler describes the endpoint callbacks for remote sync.
type Handler interface {
	Login(*pb.RsAuthenticationRequest) (*pb.RsAuthenticationResponse, error)
	Read(*pb.RsReadRequest) (*pb.RsReadResponse, error)
	Write(*pb.RsWriteRequest) (*messages.Ack, error)
	GetLastModified(*pb.RsReadRequest) (*pb.RsTimestampResponse, error)
	GetLastWrite(*pb.RsLastWriteRequest) (*pb.RsTimestampResponse, error)
	ReadDir(*pb.RsReadRequest) (*pb.RsReadDirResponse, error)
}

// StartRemoteSync starts a new RemoteSync server on the address:port specified by localServer
// and a callback interface for remote sync operations
// with given path to public and private key for TLS connection.
func StartRemoteSync(id *id.ID, localServer string, handler Handler,
	certPem, keyPem []byte) *Comms {

	// Initialize the low-level comms listeners
	pc, err := connect.StartCommServer(id, localServer,
		certPem, keyPem, nil)
	if err != nil {
		jww.FATAL.Panicf("Unable to StartCommServer: %+v", err)
	}
	rsServer := Comms{
		handler:    handler,
		ProtoComms: pc,
	}

	// Register the high-level comms endpoint functionality
	grpcServer := rsServer.GetServer()
	pb.RegisterRemoteSyncServer(grpcServer, &rsServer)
	messages.RegisterGenericServer(grpcServer, &rsServer)

	pc.ServeWithWeb()
	return &rsServer
}

// implementationFunctions for the Handler interface.
type implementationFunctions struct {
	Login           func(req *pb.RsAuthenticationRequest) (*pb.RsAuthenticationResponse, error)
	Read            func(*pb.RsReadRequest) (*pb.RsReadResponse, error)
	Write           func(*pb.RsWriteRequest) (*messages.Ack, error)
	GetLastModified func(*pb.RsReadRequest) (*pb.RsTimestampResponse, error)
	GetLastWrite    func(*messages.Ack) (*pb.RsTimestampResponse, error)
	ReadDir         func(*pb.RsReadRequest) (*pb.RsReadDirResponse, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions.
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation creates and returns a new Handler interface for implementing endpoint callbacks.
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			Login: func(*pb.RsAuthenticationRequest) (*pb.RsAuthenticationResponse, error) {
				warn(um)
				return new(pb.RsAuthenticationResponse), nil
			},
			Read: func(*pb.RsReadRequest) (*pb.RsReadResponse, error) {
				warn(um)
				return new(pb.RsReadResponse), nil
			},
			Write: func(*pb.RsWriteRequest) (*messages.Ack, error) {
				warn(um)
				return new(messages.Ack), nil
			},
			GetLastModified: func(*pb.RsReadRequest) (*pb.RsTimestampResponse, error) {
				warn(um)
				return new(pb.RsTimestampResponse), nil
			},
			GetLastWrite: func(*messages.Ack) (*pb.RsTimestampResponse, error) {
				warn(um)
				return new(pb.RsTimestampResponse), nil
			},
			ReadDir: func(*pb.RsReadRequest) (*pb.RsReadDirResponse, error) {
				warn(um)
				return new(pb.RsReadDirResponse), nil
			},
		},
	}
}

func (s *Implementation) Login(message *pb.RsAuthenticationRequest) (*pb.RsAuthenticationResponse, error) {
	return s.Functions.Login(message)
}

func (s *Implementation) Read(message *pb.RsReadRequest) (*pb.RsReadResponse, error) {
	return s.Functions.Read(message)
}
func (s *Implementation) Write(message *pb.RsWriteRequest) (*messages.Ack, error) {
	return s.Functions.Write(message)
}
func (s *Implementation) GetLastModified(message *pb.RsReadRequest) (*pb.RsTimestampResponse, error) {
	return s.Functions.GetLastModified(message)
}
func (s *Implementation) GetLastWrite(message *messages.Ack) (*pb.RsTimestampResponse, error) {
	return s.Functions.GetLastWrite(message)
}
func (s *Implementation) ReadDir(message *pb.RsReadRequest) (*pb.RsReadDirResponse, error) {
	return s.Functions.ReadDir(message)
}
