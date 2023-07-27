////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// SendAuthorizerCertRequest sends a request for an https certificate to the authorizer
func (g *Comms) SendAuthorizerCertRequest(host *connect.Host, msg *pb.AuthorizerCertRequest) (*messages.Ack, error) {
	// Create the send function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resp, err := pb.NewAuthorizerClient(conn.GetGrpcConn()).
			RequestCert(ctx, msg)
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resp)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending certificate request to authorizer: %s", msg)
	resp, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resp, result)
}

// SendEABCredentialRequest sends a request for ACME authorization to the authorizer
func (g *Comms) SendEABCredentialRequest(host *connect.Host, msg *pb.EABCredentialRequest) (*pb.EABCredentialResponse, error) {
	// Create the send function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resp, err := pb.NewAuthorizerClient(conn.GetGrpcConn()).
			RequestEABCredentials(ctx, msg)
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resp)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending certificate request to authorizer: %s", msg)
	resp, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.EABCredentialResponse{}
	return result, ptypes.UnmarshalAny(resp, result)
}
