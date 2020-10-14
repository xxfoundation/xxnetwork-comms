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
	pb "gitlab.com/elixxir/comms/mixmessages"
)

func (r *Comms) RegisterUser(registration *pb.UDBUserRegistration) *pb.UserRegistrationResponse {
	return r.handler.RegisterUser(registration)
}

func (r *Comms) RegisterFact(request *pb.FactRegisterRequest) *pb.FactRegisterResponse {
	return r.handler.RegisterFact(request)
}

func (r *Comms) ConfirmFact(request *pb.FactConfirmRequest) *pb.FactConfirmResponse {
	return r.handler.ConfirmFact(request)
}

func (r *Comms) RemoveFact(request *pb.FactRemovalRequest) *pb.FactRemovalResponse {
	return r.handler.RemoveFact(request)
}
