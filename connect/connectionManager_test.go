////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"google.golang.org/grpc"
	"math"
	"net"
	"os"
	"testing"
)

const SERVER_ADDRESS = "0.0.0.0:5556"
const SERVER_ADDRESS2 = "0.0.0.0:5557"

func TestMain(m *testing.M) {
	lis1, _ := net.Listen("tcp", SERVER_ADDRESS)
	lis2, _ := net.Listen("tcp", SERVER_ADDRESS2)

	grpcServer1 := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.MaxRecvMsgSize(33554432))

	grpcServer2 := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.MaxRecvMsgSize(33554432))

	go func() {
		_ = grpcServer1.Serve(lis1)
	}()

	go func() {
		_ = grpcServer2.Serve(lis2)
	}()
	os.Exit(m.Run())
}

// Function to test the Disconnect
// Checks if conn established in Connect() is deleted.
func TestConnectionManager_Disconnect(t *testing.T) {

	test := 2
	pass := 0
	address := SERVER_ADDRESS
	id := "pear"
	var manager ConnectionManager

	err := manager.connect(&ConnectionInfo{
		Id:             id,
		Address:        address,
		Cert:           nil,
		DisableTimeout: false,
	})
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap := manager.connections[id]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	manager.Disconnect(id)

	_, present := manager.connections[address]

	if present {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	println("Connection Manager Test: ", pass, "out of", test, "tests passed.")
}

// Function to test the Disconnect
// Checks if conn established in Connect() is deleted.
func TestConnectionManager_DisconnectAll(t *testing.T) {

	test := 4
	pass := 0
	address := SERVER_ADDRESS
	address2 := SERVER_ADDRESS2
	id := "pear"
	id2 := "apple"
	var manager ConnectionManager

	err := manager.connect(&ConnectionInfo{
		Id:             id,
		Address:        address,
		Cert:           nil,
		DisableTimeout: false,
	})
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap := manager.connections[id]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	err = manager.connect(&ConnectionInfo{
		Id:             id2,
		Address:        address2,
		Cert:           nil,
		DisableTimeout: false,
	})
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap = manager.connections[id2]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	manager.DisconnectAll()

	_, present := manager.connections[address]

	if present {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	_, present = manager.connections[address2]

	if present {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	println("Connection Manager Test: ", pass, "out of", test, "tests passed.")
}
