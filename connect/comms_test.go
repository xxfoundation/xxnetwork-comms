////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"testing"
)

// Test that trying to send to a host with no address fails
func TestSendNoAddressFails(t *testing.T) {
	// Define a new protocomms object
	comms := &ProtoComms{Id: id.NewIdFromString("test", id.Generic, t)}

	// Make fake host
	host := Host{}

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		t.Errorf("Client send function shouldn't have run")
		return nil, errors.New("Client send function shouldn't have run")
	}

	// Try to send to it and check error is right
	_, err := comms.Send(&host, f)
	if err.Error() != "Host address is blank, host might be receive only." {
		t.Errorf("Send function should have errored with address error.")
	}
}
