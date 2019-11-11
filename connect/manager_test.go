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

func TestManager_SetPrivateKey_Invalid(t *testing.T) {
	var manager Manager
	err := manager.SetPrivateKey(make([]byte, 0))
	if err == nil {
		t.Errorf("Expected error!")
	}
}

// Tests the case of obtaining a dead connection
func TestManager_ObtainConnection_DeadConnection(t *testing.T) {
	address := SERVER_ADDRESS
	var manager Manager
	host := &Host{
		address:        address,
		certificate:    nil,
		disableTimeout: false,
	}

	err := manager.AddHost(host)
	if err != nil {
		t.Errorf("Unable to create connection: %+v", err)
	}
	err = manager.connections[host.GetId()].Connection.Close()
	if err != nil {
		t.Errorf("Unable to close connection: %+v", err)
	}
	conn, err := manager.ObtainConnection(host)
	if err != nil {
		t.Errorf("Unable to obtain connection: %+v", err)
	}
	if !conn.isAlive() {
		t.Errorf("connection was not reestablished! %+v", conn)
	}
}

func TestSetCredentials_InvalidCert(t *testing.T) {
	conn := connection{
		Address:      "",
		Connection:   nil,
		Creds:        nil,
		RsaPublicKey: nil,
	}
	err := conn.setCredentials(&Host{
		address:        "",
		certificate:    []byte("test"),
		disableTimeout: false,
	})
	if err == nil {
		t.Errorf("Expected error")
	}
}

// Function to test the Disconnect
// Checks if conn established in Connect() is deleted.
func TestConnectionManager_Disconnect(t *testing.T) {

	test := 2
	pass := 0
	address := SERVER_ADDRESS
	var manager Manager
	host := &Host{
		address:        address,
		certificate:    nil,
		disableTimeout: false,
	}

	err := manager.AddHost(host)
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap := manager.connections[host.GetId()]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	manager.Disconnect(host.GetId())

	_, present := manager.connections[address]

	if present {
		t.Errorf("Disconnect Function not working properly")
	} else {
		pass++
	}

	println("connection Manager Test: ", pass, "out of", test, "tests passed.")
}

// Function to test the Disconnect
// Checks if conn established in Connect() is deleted.
func TestConnectionManager_DisconnectAll(t *testing.T) {

	test := 4
	pass := 0
	address := SERVER_ADDRESS
	address2 := SERVER_ADDRESS2
	var manager Manager
	host := &Host{
		address:        address,
		certificate:    nil,
		disableTimeout: false,
	}
	err := manager.AddHost(host)
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap := manager.connections[host.GetId()]

	if !inMap {
		t.Errorf("Connect Function didn't add connection to map")
	} else {
		pass++
	}

	err = manager.AddHost(&Host{
		address:        address2,
		certificate:    nil,
		disableTimeout: true,
	})
	if err != nil {
		t.Errorf("Unable to call connnect: %+v", err)
	}

	_, inMap = manager.connections[host.GetId()]

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

	println("connection Manager Test: ", pass, "out of", test, "tests passed.")
}

func TestConnectionManager_String(t *testing.T) {
	cm := &Manager{connections: make(map[string]*connection)}
	t.Log(cm)
	cm.connections["infoNil"] = nil
	t.Log(cm)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	// Initialize the connection object
	conn := &connection{}
	host := &Host{
		address:        "420",
		certificate:    certData,
		disableTimeout: true,
	}
	err := conn.setCredentials(host)
	if err != nil {
		t.Errorf(err.Error())
	}
	cm.connections[host.GetId()] = conn
	t.Log(cm)
}

func TestManager_SetMaxRetries(t *testing.T) {
	start := int64(10)
	cm := &Manager{
		maxRetries: start,
	}
	expected := int64(0)
	cm.SetMaxRetries(expected)
	if cm.maxRetries != expected {
		t.Errorf("Max retries did not match, got %d expected %d",
			cm.maxRetries, expected)
	}
}
