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

func TestMain(m *testing.M) {
	lis, _ := net.Listen("tcp", ":5556")

	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.MaxRecvMsgSize(33554432))

	go func() {
		defer func() { _ = lis.Close() }()
		_ = grpcServer.Serve(lis)
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

	manager.connect(id, &ConnectionInfo{
		Address: address,
	})

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

func TestConnectionManager_String(t *testing.T) {
	cm := &ConnectionManager{connections: make(map[string]*ConnectionInfo)}
	t.Log(cm)
	cm.connections["infoNil"] = nil
	t.Log(cm)
	cm.connections["fieldsNil"] = &ConnectionInfo{
		Address: "fake address",
	}
	t.Log(cm)
	// A mocked connection created without the grpc factory methods will cause
	// a panic, but there's no way to check if the field grpc uses isn't nil,
	// or to set that field up, because it's not exported
	/* cm.connections["incorrectlyCreatedConnection"] = &ConnectionInfo{
		Address: "real address",
		Connection: &grpc.ClientConn{},
	} */
}
