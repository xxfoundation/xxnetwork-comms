///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functionality for managing connections

package connect

import (
	"bytes"
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"sync"
)

// The Manager object provides thread-safe access
// to Host objects for top-level libraries
type Manager struct {
	// A map of id.IDs to Hosts
	connections sync.Map
}

// Fetch a Host from the internal map
func (m *Manager) GetHost(hostId *id.ID) (*Host, bool) {
	value, ok := m.connections.Load(*hostId)
	if !ok {
		return nil, false
	}
	host, ok := value.(*Host)
	return host, ok
}

// Creates and adds a Host object to the Manager using the given id
func (m *Manager) AddHost(id *id.ID, address string,
	cert []byte, disableTimeout, enableAuth bool) (host *Host, err error) {

	host, err = NewHost(id, address, cert, disableTimeout, enableAuth)
	if err != nil {
		return nil, err
	}

	m.addHost(host)
	return
}

// Removes a host from the connection manager
func (m *Manager) RemoveHost(id *id.ID) {
	jww.DEBUG.Printf("Removing host: %v", id)
	m.connections.Delete(*id)
}

// Internal helper function that can add Hosts directly
func (m *Manager) addHost(host *Host) {
	jww.DEBUG.Printf("Adding host: %s", host)
	m.connections.Store(*host.id, host)
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
		k := key.(id.ID)
		result.WriteString(fmt.Sprintf("[%s]: %+v",
			(&k).String(), value.(*Host)))
		return true
	})

	return result.String()
}
