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

const SERVER_ADDRESS = "localhost:5556"
const SERVER_ADDRESS2 = "localhost:5557"

func TestMain(m *testing.M) {
	lis1, _ := net.Listen("tcp", ":5556")
	lis2, _ := net.Listen("tcp", ":5557")

	grpcServer1 := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.MaxRecvMsgSize(33554432))

	grpcServer2 := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.MaxRecvMsgSize(33554432))

	go func() {
		defer func() { _ = lis1.Close() }()
		_ = grpcServer1.Serve(lis1)
	}()

	go func() {
		defer func() { _ = lis2.Close() }()
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

	manager.connect(id, address, nil, nil, false)

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

	manager.connect(id, address, nil, nil, false)

	_, inMap := manager.connections[id]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	manager.connect(id2, address2, nil, nil, false)

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

func TestConnectionManager_String(t *testing.T) {
	cm := &ConnectionManager{connections: make(map[string]*ConnectionInfo)}
	t.Log(cm)
	cm.connections["infoNil"] = nil
	t.Log(cm)
	cm.connections["fieldsNil"] = &ConnectionInfo{
		Address: "fake address",
	}
	t.Log(cm)
	// A mocked connection created without the gRPC factory methods will cause
	// a panic, but there's no way to check if the field gRPC uses isn't nil,
	// or to set that field up, because it's not exported
	/* cm.connections["incorrectlyCreatedConnection"] = &ConnectionInfo{
		Address: "real address",
		Connection: &grpc.ClientConn{},
	} */
}
