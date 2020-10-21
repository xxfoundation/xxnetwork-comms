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
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/messages"
)

func (r *Comms) RegisterUser(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	registration := &pb.UDBUserRegistration{}
	err = ptypes.UnmarshalAny(msg.Message, registration)
	if err != nil {
		return nil, err
	}

	return r.handler.RegisterUser(registration, authState)
}

func (r *Comms) RegisterFact(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.FactRegisterResponse, error) {
	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	request := &pb.FactRegisterRequest{}
	err = ptypes.UnmarshalAny(msg.Message, request)
	if err != nil {
		return nil, err
	}

	return r.handler.RegisterFact(request, authState)
}

func (r *Comms) ConfirmFact(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.Fact, error) {
	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	request := &pb.FactConfirmRequest{}
	err = ptypes.UnmarshalAny(msg.Message, request)
	if err != nil {
		return nil, err
	}

	return r.handler.ConfirmFact(request, authState)
}

func (r *Comms) RemoveFact(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	// Create an auth object
	authState, err := r.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unmarshall the any message to the message type needed
	request := &pb.FactRemovalRequest{}
	err = ptypes.UnmarshalAny(msg.Message, request)
	if err != nil {
		return nil, err
	}

	return r.handler.RemoveFact(request, authState)
}
