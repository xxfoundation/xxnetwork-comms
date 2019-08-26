////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server gRPC endpoints

package registration

import (
	"fmt"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
)

// RegisterUser event handler which registers a user with the platform
func (r *RegistrationComms) RegisterUser(ctx context.Context, msg *pb.UserRegistration) (
	*pb.UserRegistrationConfirmation, error) {
	// Obtain the signed key by passing to registration server
	pubKey := msg.GetClientRSAPubKey()
	signature, err := r.handler.RegisterUser(msg.GetRegistrationCode(), pubKey)
	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
		err = errors.New(err.Error())
	}

	// Return the confirmation message
	return &pb.UserRegistrationConfirmation{
		ClientSignedByServer: &pb.RSASignature{
			Signature: signature,
		},
		Error: errMsg,
	}, err
}

// CheckClientVersion event handler which checks whether the client library
// version is compatible with the network
func (r *RegistrationComms) CheckClientVersion(ctx context.Context, msg *pb.ClientVersion) (*pb.ClientVersionConfirmation, error) {
	isOK, err := r.handler.CheckClientVersion(msg.Version)

	// Return the confirmation message
	return &pb.ClientVersionConfirmation{
		IsOK:isOK,
	}, err
}

// Handle a node registration event
func (r *RegistrationComms) RegisterNode(ctx context.Context, msg *pb.NodeRegistration) (
	*pb.Ack, error) {
	// Obtain peer IP address
	info, _ := peer.FromContext(ctx)
	host, _, err := net.SplitHostPort(info.Addr.String())
	if err != nil {
		return &pb.Ack{}, err
	}
	addr := fmt.Sprintf("%s:%s", host, msg.GetPort())

	// Pass information for Node registration
	err = r.handler.RegisterNode(msg.GetID(), addr, msg.GetServerTlsCert(),
		msg.GetGatewayAddress(), msg.GetGatewayTlsCert(),
		msg.GetRegistrationCode())
	return &pb.Ack{}, err
}
