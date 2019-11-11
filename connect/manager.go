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
	"math"
	"sync"
)

// The Manager object provides thread-safe access
// to Host objects for top-level libraries
type Manager struct {
	// A map of string IDs to Hosts
	connections sync.Map
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

// Initializes a host object and adds the newly-created object to the Manager
func (m *Manager) AddHost(id, address string, cert []byte,
	disableTimeout bool) (err error) {

	// Initialize the Host object
	host := &Host{
		address:     address,
		certificate: cert,
	}

	// Set the max number of retries for establishing a connection
	if disableTimeout {
		host.maxRetries = math.MaxInt64
	} else {
		host.maxRetries = 100
	}

	// Configure the host credentials
	err = host.setCredentials()
	if err != nil {
		return
	}

	// Add the connection to the manager
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
