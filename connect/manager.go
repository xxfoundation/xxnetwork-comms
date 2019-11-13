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

// Adds the given Host object to the Manager using the given id
func (m *Manager) AddHost(id string, host *Host) {
	m.connections.Store(id, host)
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
