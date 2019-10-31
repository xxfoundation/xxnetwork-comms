package connect

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"golang.org/x/net/context"
	"math"
	"sort"
	"sync"
	"time"
)

// Default maximum number of retries
const DefaultMaxRetries = 100

type Manager struct {
	// A map of string IDs to open connections
	connections map[string]*connection

	// Local server RSA Private Key
	privateKey *rsa.PrivateKey

	lock       sync.Mutex
	maxRetries int64
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

// Gets a connection object from the Manager
// Or creates and returns a new one if it does not already exist
func (m *Manager) ObtainConnection(connInfo *Host) (*connection, error) {
	// If the connection does not already exist, create a new connection
	conn, ok := m.connections[connInfo.Id]
	if !ok {
		jww.INFO.Printf("Connection %s does not exist, creating...",
			connInfo.Id)
		err := m.createConnection(connInfo)
		if err != nil {
			return nil, err
		}
		conn, ok = m.connections[connInfo.Id]
	}

	// Verify the connection is still good
	if !conn.isAlive() {
		// If not, attempt to reestablish the connection
		jww.WARN.Printf("Bad connection state, reconnecting: %v",
			conn.Connection)
		m.Disconnect(connInfo.Id)
		err := m.createConnection(connInfo)
		if err != nil {
			return nil, err
		}
		conn, ok = m.connections[connInfo.Id]
	}
	return conn, nil
}

// Set the global max number of retries for establishing connections
func (m *Manager) SetMaxRetries(mr int64) {
	m.maxRetries = mr
}

// Fully initializes a connection object given a Host object
// Then adds the newly-created object to the Manager connection map
func (m *Manager) createConnection(host *Host) (err error) {
	// Initialize the connections map if it hasn't already
	if m.connections == nil {
		m.connections = make(map[string]*connection)
	}

	// Initialize the connection object
	conn := &connection{
		Address: host.Address,
	}

	// Obtain credentials for the connection object
	err = conn.setCredentials(host)
	if err != nil {
		return
	}

	// Set the max number of retries for establishing a connection
	var maxRetries int64
	if host.DisableTimeout {
		maxRetries = math.MaxInt64
	} else {
		if m.maxRetries == 0 {
			maxRetries = DefaultMaxRetries
		} else {
			maxRetries = m.maxRetries
		}
	}

	// Establish the connection
	err = conn.connect(maxRetries)

	// Add the connection to the manager
	m.lock.Lock()
	m.connections[host.Id] = conn
	m.lock.Unlock()
	return
}

// Closes a client connection and removes it from Manager
func (m *Manager) Disconnect(id string) {
	m.lock.Lock()
	connection, present := m.connections[id]
	if present {
		err := connection.Connection.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %+v", id,
				errors.New(err.Error()))
		}
		delete(m.connections, id)
	}
	m.lock.Unlock()
}

// Closes all client connections and removes them from Manager
func (m *Manager) DisconnectAll() {
	for connId := range m.connections {
		m.Disconnect(connId)
	}
}

// Implements Stringer for debug printing
func (m *Manager) String() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Sort connection IDs to print in a consistent order
	keys := make([]string, len(m.connections))
	i := 0
	for key := range m.connections {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// Print each connection's information
	var result bytes.Buffer
	for _, key := range keys {
		// Populate fields without ever dereferencing nil
		connection := m.connections[key]
		if connection != nil {
			result.WriteString(fmt.Sprintf(
				"[%v] %s",
				key, connection))
		}
	}

	return result.String()
}

// TimeoutContext is a context with the default timeout
func ConnectionContext(seconds time.Duration) (context.Context, context.CancelFunc) {
	waitingPeriod := seconds * time.Second
	jww.DEBUG.Printf("Timing out in: %s", waitingPeriod)
	ctx, cancel := context.WithTimeout(context.Background(),
		waitingPeriod)
	return ctx, cancel

}

// DefaultContexts creates a context object with the default context
// for all client messages. This is primarily used to set the default
// timeout for all clients
func MessagingContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(),
		2*time.Minute)
	return ctx, cancel
}

// StreamingContext creates a context object with the default context
// for all client streaming messages. This is primarily used to
// allow a cancel option for clients and is suitable for unary streaming.
func StreamingContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel
}
