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
	"gitlab.com/elixxir/primitives/id"
	"sync"
	"testing"
)

// The Manager object provides thread-safe access
// to Host objects for top-level libraries
type Manager struct {
	// A map of id.IDs to Hosts
	connections map[id.ID]*Host
	mux         sync.RWMutex
}

func newManager() *Manager {
	return &Manager{
		connections: make(map[id.ID]*Host),
		mux:         sync.RWMutex{},
	}
}

func NewManagerTesting(t *testing.T) *Manager {
	if t == nil {
		jww.FATAL.Panicf("NewMangerTesting is for testing only")
	}
	return newManager()
}

// Fetch a Host from the internal map
func (m *Manager) GetHost(hostId *id.ID) (*Host, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	host, ok := m.connections[*hostId]
	if !ok {
		return nil, false
	}
	return host, ok
}

// Creates and adds a Host object to the Manager using the given id
func (m *Manager) AddHost(hid *id.ID, address string,
	cert []byte, params HostParams) (host *Host, err error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	//check if the host already exists, if it does return it
	host, ok := m.connections[*hid]
	if ok {
		return host, nil
	}

	//create the new host
	host, err = NewHost(hid, address, cert, params)
	if err != nil {
		return nil, err
	}

	//add the host to the map
	m.addHost(host)

	return host, nil
}

func (m *Manager) addHost(host *Host) {
	jww.DEBUG.Printf("Adding host: %s", host)
	m.connections[*(host.id)] = host
}

// Removes a host from the connection manager
func (m *Manager) RemoveHost(hid *id.ID) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.connections, *hid)
}

// Closes all client connections and removes them from Manager
func (m *Manager) DisconnectAll() {
	m.mux.RLock()
	defer m.mux.RUnlock()
	for _, host := range m.connections {
		host.Disconnect()
	}
}

// Implements Stringer for debug printing
func (m *Manager) String() string {
	var result bytes.Buffer
	for k, host := range m.connections {
		result.WriteString(fmt.Sprintf("[%s]: %+v",
			(&k).String(), host))
	}
	return result.String()
}
