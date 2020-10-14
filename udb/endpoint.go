///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains user discovery server gRPC endpoint wrappers
// When you add the udb server to mixmessages/mixmessages.proto and add the
// first function, a version of that goes here which calls the "handler"
// version of the function, with any mappings/wrappings necessary.

package udb

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// ClientCall implements a call from client -> UDB
// msg is here as a dummy, it's a parameter that would be defined in
// mixmessages.proto against this ClientCall function. Same thing with the
// return call. The parameters here have to line up to what is seen
// in mixmessages.proto.
// In this specific case, that would look like:
//
// rpc ClientCall (messages.AuthenticatedMessage) returns (ServerPollResponse){}
func (r *Comms) ClientCall(ctx context.Context,
	msg *messages.AuthenticatedMessage) (
	*pb.PermissionPollResponse, error) {

	// If we were actually using authenticated messages, we would need
	// to validate and unwrap it, as follows:

	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf(
			"Unable handles reception of AuthenticatedMessage: %+v",
			err)
	}

	// Unmarshall the any message to the message type needed
	pollMsg := &pb.PermissioningPoll{}
	err = ptypes.UnmarshalAny(msg.Message, pollMsg)
	if err != nil {
		return nil, err
	}

	// Now that we have unmarshalled the message we do any "extra stuff" we
	// would need. Again it is unlikely you ever need to do this, but you
	// can in special circumstances. This extracts the connecting ip from
	// the connection object.

	// Get server IP and port
	ip, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return &pb.PermissionPollResponse{}, err
	}
	port := pollMsg.ServerPort
	address := fmt.Sprintf("%s:%d", ip, port)

	// Now we call the "ClientCall" function defined in handler.go, which
	// will return that it is unimplemented.
	// In 99% of cases, this is the only line you will need in the
	// endpoint.go function:
	return r.handler.ClientCall(pollMsg, authState, address)
}

func (r *Comms) RegisterUser(registration *pb.UDBUserRegistration) pb.UserRegistrationResponse {
	return r.handler.RegisterUser(registration)
}

func (r *Comms) RegisterFact(request *pb.FactRegisterRequest) pb.FactRegisterResponse {
	return r.handler.RegisterFact(request)
}

func (r *Comms) ConfirmFact(request *pb.FactConfirmRequest) pb.FactConfirmResponse {
	return r.handler.ConfirmFact(request)
}

func (r *Comms) RemoveFact(request *pb.FactRemovalRequest) pb.FactRemovalResponse {
	return r.handler.RemoveFact(request)
}
