///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains registration server gRPC endpoints

package registration

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (r *Comms) AuthenticateToken(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	err := r.ValidateToken(msg)
	if err != nil {
		jww.ERROR.Printf("Unable to authenticate token: %+v", err)
	}
	return &messages.Ack{}, err
}

// Handles reception of reverse-authentication token requests
func (r *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := r.GenerateToken()
	return &messages.AssignToken{
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
		ClientSignedByServer: &messages.RSASignature{
			Signature: signature,
		},
		Error: errMsg,
	}, err
}

// CheckClientVersion event handler which checks whether the client library
// version is compatible with the network
func (r *Comms) GetCurrentClientVersion(ctx context.Context, ping *messages.Ping) (*pb.ClientVersion, error) {

	version, err := r.handler.GetCurrentClientVersion()

	// Return the confirmation message
	return &pb.ClientVersion{
		Version: version,
	}, err
}

// Handle a node registration event
func (r *Comms) RegisterNode(ctx context.Context, msg *pb.NodeRegistration) (
	*messages.Ack, error) {

	// Infer peer IP address (do not use msg.GetServerAddress())
	ip, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return &messages.Ack{}, err
	}

	port := msg.GetServerPort()
	address := fmt.Sprintf("%s:%d", ip, port)

	gwAddress := fmt.Sprintf("%s:%d", msg.GetGatewayAddress(),
		msg.GetGatewayPort())

	// Pass information for Node registration
	err = r.handler.RegisterNode(msg.GetSalt(), address,
		msg.GetServerTlsCert(),
		gwAddress, msg.GetGatewayTlsCert(),
		msg.GetRegistrationCode())
	return &messages.Ack{}, err
}

// Handles incoming requests for the NDF
func (r *Comms) PollNdf(ctx context.Context, ndfHash *pb.NDFHash) (*pb.NDF, error) {
	newNDF, err := r.handler.PollNdf(ndfHash.Hash)
	//Return the new ndf
	return &pb.NDF{Ndf: newNDF}, err
}

// Server -> Permissioning unified polling
func (r *Comms) Poll(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.PermissionPollResponse, error) {
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

// Server -> Permissioning unified polling
func (r *Comms) CheckRegistration(ctx context.Context, msg *pb.RegisteredNodeCheck) (*pb.RegisteredNodeConfirmation, error) {

	//Return the new ndf
	return r.handler.CheckRegistration(msg)
}
