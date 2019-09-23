////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"testing"
)

// Smoke test SendRegistrationMessage
func TestSendRegistrationMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	rg := registration.StartRegistrationServer(GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	connID := MockID("clientToRegistration")
	var c ClientComms
	c.ConnectToRemote(connID, GatewayAddress, nil, false)

	_, err := c.SendRegistrationMessage(connID, &pb.UserRegistration{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckClientVersion
func TestSendCheckClientVersionMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	rg := registration.StartRegistrationServer(GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	connID := MockID("clientToRegistration")
	var c ClientComms
	c.ConnectToRemote(connID, GatewayAddress, nil, false)

	_, err := c.SendGetCurrentClientVersionMessage(connID)
	if err != nil {
		t.Errorf("CheckClientVersion: Error received: %s", err)
	}
}

//Call GetUpdatedNDF on the registration server
func (c *ClientComms) SendGetUpdatedNDF(id fmt.Stringer, message *pb.NDFHash) (*pb.NDF, error) {
	//Get the connection
	connection := c.GetRegistrationConnection(id)
	ctx, cancel := connect.MessagingContext()

	//Send message
	response, err := connection.GetUpdatedNDF(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("GetUpdatedNDf: Error received: %v", err)
	}

	cancel()
	return response, err
}
