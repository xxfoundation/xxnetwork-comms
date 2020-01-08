////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"gitlab.com/elixxir/comms/testkeys"
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

func TestSetCredentials_InvalidCert(t *testing.T) {
	host := &Host{
		address:     "",
		certificate: []byte("test"),
	}
	err := host.setCredentials()
	if err == nil {
		t.Errorf("Expected error")
	}
}

// Function to test the Disconnect
// Checks if conn established in connect() is deleted.
func TestConnectionManager_Disconnect(t *testing.T) {

	test := 2
	pass := 0
	address := SERVER_ADDRESS
	var manager Manager
	testId := "testId"
	host, err := manager.AddHost(testId, address, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap := manager.connections.Load(testId)

	if !inMap {
		t.Errorf("connect Function didn't add connection to map")
	} else {
		pass++
	}

	f := func(host *Host) error {
		return nil
	}

	err = host.connect(f)
	if err != nil {
		t.Error("Unable to connect")
	}
	host.Disconnect()

	if host.isAlive() {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	println("connection Manager Test: ", pass, "out of", test, "tests passed.")
}

// Function to test the Disconnect
// Checks if conn established in connect() is deleted.
func TestConnectionManager_DisconnectAll(t *testing.T) {

	test := 4
	pass := 0
	address := SERVER_ADDRESS
	address2 := SERVER_ADDRESS2
	var manager Manager
	testId := "testId"
	testId2 := "TestId2"

	host, err := manager.AddHost(testId, address, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, inMap := manager.GetHost(testId)

	if !inMap {
		t.Errorf("connect Function didn't add connection to map")
	} else {
		pass++
	}

	host2, err := manager.AddHost(testId2, address2, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	f := func(host *Host) error {
		return nil
	}

	err = host.connect(f)
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}
	err = host2.connect(f)
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap = manager.connections.Load(testId2)

	if !inMap {
		t.Errorf("connect Function didn't add connection to map")
	} else {
		pass++
	}

	manager.DisconnectAll()

	if host.isAlive() {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}
	if host2.isAlive() {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	println("connection Manager Test: ", pass, "out of", test, "tests passed.")
}

func TestConnectionManager_String(t *testing.T) {
	var manager Manager
	t.Log(manager)

	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	_, err := manager.AddHost("test", "test", certData, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Initialize the connection object
	t.Log(manager.String())
}
