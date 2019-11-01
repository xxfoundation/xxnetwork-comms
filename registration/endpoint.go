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
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// RegisterUser event handler which registers a user with the platform
func (r *Comms) RegisterUser(ctx context.Context, msg *pb.UserRegistration) (
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
func (r *Comms) GetCurrentClientVersion(ctx context.Context, msg *pb.Ping) (*pb.ClientVersion, error) {
	version, err := r.handler.GetCurrentClientVersion()

	// Return the confirmation message
	return &pb.ClientVersion{
		Version: version,
	}, err
}

// Handle a node registration event
func (r *Comms) RegisterNode(ctx context.Context, msg *pb.NodeRegistration) (
	*pb.Ack, error) {
	// Obtain peer IP address
	host, _, err := connect.GetAddressFromContext(ctx)
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

//GetUpdatedNDF event handler handles a client's request for a new ndf on the permissioning server
func (r *Comms) GetUpdatedNDF(ctx context.Context, msg *pb.NDFHash) (*pb.NDF, error) {
	newNDF, err := r.handler.GetUpdatedNDF(msg.Hash)
	//Return the new ndf
	return &pb.NDF{Ndf: newNDF}, err
}
