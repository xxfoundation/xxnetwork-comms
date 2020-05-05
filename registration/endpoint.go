////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server gRPC endpoints

package registration

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (r *Comms) AuthenticateToken(ctx context.Context,
	msg *pb.AuthenticatedMessage) (*pb.Ack, error) {
	err := r.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &pb.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (r *Comms) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	token, err := r.GenerateToken()
	return &pb.AssignToken{
		Token: token,
	}, err
}

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
func (r *Comms) GetCurrentClientVersion(ctx context.Context, ping *pb.Ping) (*pb.ClientVersion, error) {

	version, err := r.handler.GetCurrentClientVersion()

	// Return the confirmation message
	return &pb.ClientVersion{
		Version: version,
	}, err
}

// Handle a node registration event
func (r *Comms) RegisterNode(ctx context.Context, msg *pb.NodeRegistration) (
	*pb.Ack, error) {

	// Infer peer IP address (do not use msg.GetServerAddress())
	ip, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return &pb.Ack{}, err
	}

	port := msg.GetServerPort()
	address := fmt.Sprintf("%s:%d", ip, port)

	gwAddress := fmt.Sprintf("%s:%d", msg.GetGatewayAddress(),
		msg.GetGatewayPort())

	// Pass information for Node registration
	err = r.handler.RegisterNode(msg.GetID(), address,
		msg.GetServerTlsCert(),
		gwAddress, msg.GetGatewayTlsCert(),
		msg.GetRegistrationCode())
	return &pb.Ack{}, err
}

// Handles incoming requests for the NDF
func (r *Comms) PollNdf(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.NDF, error) {
	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	ndfHash := &pb.NDFHash{}
	err = ptypes.UnmarshalAny(msg.Message, ndfHash)
	if err != nil {
		return nil, err
	}

	newNDF, err := r.handler.PollNdf(ndfHash.Hash, authState)
	//Return the new ndf
	return &pb.NDF{Ndf: newNDF}, err
}

// Server -> Permissioning unified polling
func (r *Comms) Poll(ctx context.Context, msg *pb.AuthenticatedMessage) (*pb.PermissionPollResponse, error) {
	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	pollMsg := &pb.PermissioningPoll{}
	err = ptypes.UnmarshalAny(msg.Message, pollMsg)
	if err != nil {
		return nil, err
	}

	//Return the new ndf
	return r.handler.Poll(pollMsg, authState)
}
