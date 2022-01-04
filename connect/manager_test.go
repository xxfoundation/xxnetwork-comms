///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"math"
	"net"
	"os"
	"testing"
)

const ServerAddress = "0.0.0.0:5556"
const ServerAddress2 = "0.0.0.0:5557"

func TestMain(m *testing.M) {
	lis1, _ := net.Listen("tcp", ServerAddress)
	lis2, _ := net.Listen("tcp", ServerAddress2)
	TestingOnlyDisableTLS = true

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
	address := ServerAddress
	manager := newManager()
	testId := id.NewIdFromString("testId", id.Node, t)
	host, err := manager.AddHost(testId, address, nil, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap := manager.connections[*testId]

	if !inMap {
		t.Errorf("connect Function didn't add connection to map")
	} else {
		pass++
	}

	err = host.connect()
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
	address := ServerAddress
	address2 := ServerAddress2
	manager := newManager()
	testId := id.NewIdFromString("testId", id.Generic, t)
	testId2 := id.NewIdFromString("TestId2", id.Generic, t)

	host, err := manager.AddHost(testId, address, nil, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, inMap := manager.GetHost(testId)

	if !inMap {
		t.Errorf("connect Function didn't add connection to map")
	} else {
		pass++
	}

	host2, err := manager.AddHost(testId2, address2, nil, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	err = host.connect()
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}
	err = host2.connect()
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap = manager.connections[*testId2]

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
	manager := newManager()
	//t.Log(manager)

	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	testID := id.NewIdFromString("test", id.Node, t)
	_, err := manager.AddHost(testID, "test", certData, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	// Initialize the connection object
	t.Log(manager.String())
}

// Show that if a connection is in the map,
// it's no longer in the map after RemoveHost is called
func TestConnectionManager_RemoveHost(t *testing.T) {
	manager := newManager()

	// After adding the host, the connection should be accessible
	id := id.NewIdFromString("i am a connection", id.Gateway, t)
	manager.addHost(&Host{id: id})
	_, ok := manager.GetHost(id)
	if !ok {
		t.Errorf("Host with id %v not in connection manager", id)
	}

	// After removing the host, the connection should no longer be accessible
	// from the manager
	manager.RemoveHost(id)
	_, ok = manager.GetHost(id)
	if ok {
		t.Errorf("Host with id %v was in connection manager, but oughtn't to have been", id)
	}

	// Removing the host again shouldn't cause any panics or problems
	manager.RemoveHost(id)
}
