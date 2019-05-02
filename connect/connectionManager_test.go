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
		defer lis.Close()
		grpcServer.Serve(lis)
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

	connect(id, &ConnectionInfo{
		Address: address,
	})

	_, inMap := connections[id]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	Disconnect(id)

	_, present := connections[address]

	if present {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	println("Connection Manager Test: ", pass, "out of", test, "tests passed.")
}
