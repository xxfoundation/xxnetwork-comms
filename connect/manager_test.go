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

// Tests the case of obtaining a dead connection
func TestManager_ObtainConnection_DeadConnection(t *testing.T) {
	address := SERVER_ADDRESS
	id := "pear"
	var manager Manager
	host := &Host{
		Id:             id,
		Address:        address,
		Cert:           nil,
		DisableTimeout: false,
	}

	err := manager.createConnection(host)
	if err != nil {
		t.Errorf("Unable to create connection: %+v", err)
	}
	err = manager.connections[id].Connection.Close()
	if err != nil {
		t.Errorf("Unable to close connection: %+v", err)
	}
	conn, err := manager.ObtainConnection(host)
	if err != nil {
		t.Errorf("Unable to obtain connection: %+v", err)
	}
	if !conn.isAlive() {
		t.Errorf("Connection was not reestablished! %+v", conn)
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
		Id:             "",
		Address:        "",
		Cert:           []byte("test"),
		DisableTimeout: false,
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
	id := "pear"
	var manager Manager

	err := manager.createConnection(&Host{
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
	var manager Manager

	err := manager.createConnection(&Host{
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

	err = manager.createConnection(&Host{
		Id:             id2,
		Address:        address2,
		Cert:           nil,
		DisableTimeout: true,
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

func TestConnectionManager_String(t *testing.T) {
	cm := &Manager{connections: make(map[string]*connection)}
	t.Log(cm)
	cm.connections["infoNil"] = nil
	t.Log(cm)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)
	id := "420"
	// Initialize the connection object
	conn := &connection{}
	err := conn.setCredentials(&Host{
		Id:             id,
		Address:        "420",
		Cert:           certData,
		DisableTimeout: true,
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	cm.connections[id] = conn
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
