///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package interconnect

import (
	"bytes"
	"context"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

func TestComms_GetNDF(t *testing.T) {

	testNodeID := id.NewIdFromString("test", id.Node, t)
	testPort := "5959"

	certPEM := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	keyPEM := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())

	ic, _ := StartCMixInterconnect(testNodeID, testPort, NewImplementation(), certPEM, keyPEM)

	expectedMessage := []byte("hello world")

	resultMsg, err := ic.GetNDF(context.Background(), &messages.Ping{})
	if err != nil {
		t.Errorf("Failed to send message: %v", err)
	}
	if !bytes.Equal(expectedMessage, resultMsg.Ndf) {
		t.Errorf("Unexpected message. "+
			"\nReceived: %v"+
			"\nExpected: %v", resultMsg.Ndf, expectedMessage)
	}

}
