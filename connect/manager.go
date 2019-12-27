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

// Creates and adds a Host object to the Manager using the given id
func (m *Manager) AddHost(id, address string,
	cert []byte, disableTimeout, enableAuth bool) (host *Host, err error) {

	host, err = NewHost(id, address, cert, disableTimeout, enableAuth)
	if err != nil {
		return nil, err
	}

	m.connections.Store(id, host)
	return
}

// Closes all client connections and removes them from Manager
func (m *Manager) DisconnectAll() {
	m.connections.Range(func(key interface{}, value interface{}) bool {
		value.(*Host).Disconnect()
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
