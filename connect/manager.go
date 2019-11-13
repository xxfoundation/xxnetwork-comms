////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for managing connections

package connect

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"sync"
)

// The Manager object provides thread-safe access
// to Host objects for top-level libraries
type Manager struct {
	// A map of string IDs to Hosts
	connections sync.Map
	// Private key of the local communication server
	privateKey *rsa.PrivateKey
}

// Set private key to data to a PEM block
func (m *Manager) SetPrivateKey(data []byte) error {
	key, err := rsa.LoadPrivateKeyFromPem(data)
	if err != nil {
		s := fmt.Sprintf("Failed to form private key file from data at %s: %+v", data, err)
		return errors.New(s)
	}

	m.privateKey = key
	return nil
}

// Get connection manager's private key
func (m *Manager) GetPrivateKey() *rsa.PrivateKey {
	return m.privateKey
}

// Fetch a Host from the internal map
func (m *Manager) GetHost(hostId string) (*Host, bool) {
	value, ok := m.connections.Load(hostId)
	if !ok {
		return nil, false
	}
	host, ok := value.(*Host)
	return host, ok
}

// Creates a host object and adds the newly-created object to the Manager
func (m *Manager) AddHost(id, address string, cert []byte,
	disableTimeout bool) (host *Host, err error) {

	// Create the Host object
	host, err = NewHost(address, cert, disableTimeout)
	if err != nil {
		return
	}

	// Add the Host to the Manager
	m.connections.Store(id, host)
	return
}

// Closes all client connections and removes them from Manager
func (m *Manager) DisconnectAll() {
	m.connections.Range(func(key interface{}, value interface{}) bool {
		value.(*Host).disconnect()
		return true
	})
}

// Implements Stringer for debug printing
func (m *Manager) String() string {
	var result bytes.Buffer
	m.connections.Range(func(key interface{}, value interface{}) bool {
		result.WriteString(fmt.Sprintf("[%s]: %+v",
			key.(string), value.(*Host)))
		return true
	})

	return result.String()
}
