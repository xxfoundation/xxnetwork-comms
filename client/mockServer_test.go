///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains a dummy/mock server instance for testing purposes

package client

import (
	"fmt"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"gitlab.com/xx_network/comms/connect"
	"sync"
)

// ------------------------- Testing globals -------------------------------------

var GetHostErrBool = true
var RequestNdfErr error = nil
var NdfToreturn = pb.NDF{Ndf: []byte(testutils.ExampleJSON)}

const RegistrationAddr = "0.0.0.0:5558"

var ExampleBadNdfJSON = "badNDF"
var RegistrationHandler = &MockRegistration{}

const RegistrationAddrErr = "0.0.0.0:5559"

var portLock sync.Mutex
var port = 5800
var RegistrationError = &MockRegistrationError{}
var Retries = 0

// Utility function to avoid address collisions in testing suite
func getNextAddress() string {
	portLock.Lock()
	defer func() {
		port++
		portLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", port)
}

// ------------------------- Mock Registration Server Handler ---------------------------

type MockRegistration struct {
}

func (s *MockRegistration) RegisterNode(salt []byte, serverAddr, serverTlsCert, gatewayAddr,
	gatewayTlsCert, registrationCode string) error {
	return nil
}

func (s *MockRegistration) PollNdf(clientNdfHash []byte, auth *connect.Auth) ([]byte, error) {
	return []byte(testutils.ExampleJSON), nil
}

func (s *MockRegistration) Poll(*pb.PermissioningPoll, *connect.Auth, string) (*pb.PermissionPollResponse, error) {
	return &pb.PermissionPollResponse{}, nil
}

// Registers a user and returns a signed public key
func (s *MockRegistration) RegisterUser(registrationCode,
	key string) (hash []byte, err error) {
	return nil, nil
}

func (s *MockRegistration) GetCurrentClientVersion() (version string, err error) {
	return "", nil
}

func (s *MockRegistration) CheckRegistration(msg *pb.RegisteredNodeCheck) (*pb.RegisteredNodeConfirmation, error) {
	return nil, nil
}

// ------------------------- Mock Error Registration Server Handler ---------------------------

type MockRegistrationError struct {
}

func (s *MockRegistrationError) RegisterNode(salt []byte, serverAddr, serverTlsCert, gatewayAddr,
	gatewayTlsCert, registrationCode string) error {
	return nil
}

func (s *MockRegistrationError) PollNdf(clientNdfHash []byte, auth *connect.Auth) ([]byte, error) {
	if Retries < 5 {
		Retries++
		return nil, errors.New(ndf.NO_NDF)
	}
	return []byte(testutils.ExampleJSON), nil
}

func (s *MockRegistrationError) Poll(*pb.PermissioningPoll, *connect.Auth, string) (*pb.PermissionPollResponse, error) {
	if Retries < 5 {
		Retries++
		return nil, errors.New(ndf.NO_NDF)
	}
	return &pb.PermissionPollResponse{}, nil
}

// Registers a user and returns a signed public key
func (s *MockRegistrationError) RegisterUser(registrationCode,
	key string) (hash []byte, err error) {
	return nil, nil
}

func (s *MockRegistrationError) GetCurrentClientVersion() (version string, err error) {
	return "", nil

}

func (s *MockRegistrationError) CheckRegistration(msg *pb.RegisteredNodeCheck) (*pb.RegisteredNodeConfirmation, error) {
	return nil, nil
}
