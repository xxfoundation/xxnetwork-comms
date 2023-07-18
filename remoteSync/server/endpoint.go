////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains remote sync gRPC endpoints

package server

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// Login to the server, receiving a token
func (rc *Comms) Login(ctx context.Context, message *pb.RsAuthenticationRequest) (*pb.RsAuthenticationResponse, error) {
	return rc.handler.Login(message)
}

// Read data from the server
func (rc *Comms) Read(ctx context.Context, message *pb.RsReadRequest) (*pb.RsReadResponse, error) {
	return rc.handler.Read(message)
}

// Write data to the server
func (rc *Comms) Write(ctx context.Context, message *pb.RsWriteRequest) (*messages.Ack, error) {
	return rc.handler.Write(message)
}

// GetLastModified returns the last time a resource was modified
func (rc *Comms) GetLastModified(ctx context.Context, message *pb.RsReadRequest) (*pb.RsTimestampResponse, error) {
	return rc.handler.GetLastModified(message)
}

// GetLastWrite returns the last time this remote sync server was modified
func (rc *Comms) GetLastWrite(ctx context.Context, message *pb.RsLastWriteRequest) (*pb.RsTimestampResponse, error) {
	return rc.handler.GetLastWrite(message)
}

// ReadDir reads a directory from the server
func (rc *Comms) ReadDir(ctx context.Context, message *pb.RsReadRequest) (*pb.RsReadDirResponse, error) {
	return rc.handler.ReadDir(message)
}
