////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Manager struct and operational functions

package gossip

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/primitives/id"
	"sync"
	"time"
)

// Structure holding messages for a given tag, if the tag does not yet exist
// If the tag is not created in 5 minutes, the record should be deleted
type MessageRecord struct {
	Timestamp time.Time
	Messages  []*GossipMsg
}

type ManagerFlags struct {
	// How long a message record should last in the buffer
	BufferExpirationTime time.Duration

	// Frequency with which to check the buffer.
	// Should be long, since the thread takes a lock each time it checks the buffer
	MonitorThreadFrequency time.Duration
}

func DefaultManagerFlags() ManagerFlags {
	return ManagerFlags{
		BufferExpirationTime:   300 * time.Second,
		MonitorThreadFrequency: 150 * time.Second,
	}
}

// Manager for various GossipProtocols that are accessed by tag
type Manager struct {
	comms *connect.ProtoComms

	// Stored map of GossipProtocols
	protocols    map[string]*Protocol
	protocolLock sync.RWMutex // Lock for protocols object

	// Buffer messages with tags that do not have a protocol created yet
	buffer     map[string]*MessageRecord // TODO: should this be sync.Map?
	bufferLock sync.RWMutex              // Lock for buffers object

	flags ManagerFlags

	*UnimplementedGossipServer
}

// Creates a new Gossip Manager struct
func NewManager(comms *connect.ProtoComms, flags ManagerFlags) *Manager {
	m := &Manager{
		comms:     comms,
		protocols: map[string]*Protocol{},
		buffer:    map[string]*MessageRecord{},
		flags:     flags,
	}
	_ = m.bufferMonitor()
	return m
}

// Creates and stores a new Protocol in the Manager
func (m *Manager) NewGossip(tag string, flags ProtocolFlags,
	receiver Receiver, verifier SignatureVerification, peers []*id.ID) {
	m.protocolLock.Lock()
	defer m.protocolLock.Unlock()

	protocol := &Protocol{
		fingerprints:    map[Fingerprint]*uint64{},
		oldFingerprints: map[Fingerprint]*uint64{},
		comms:           m.comms,
		peers:           peers,
		flags:           flags,
		receiver:        receiver,
		verify:          verifier,
		IsDefunct:       false,
		crand:           csprng.NewSystemRNG(),
		sendWorkers:     make(chan sendInstructions, 100*flags.NumParallelSends),
		fingerprinter:   flags.Fingerprinter,
	}

	// Set default fingerprinter function
	if flags.Fingerprinter == nil {
		protocol.fingerprinter = func(msg *GossipMsg) Fingerprint {
			return getFingerprint(msg)
		}
	}

	// create the runners
	launchSendWorkers(flags.NumParallelSends, protocol.sendWorkers)

	m.protocols[tag] = protocol

	m.bufferLock.Lock()
	if record, ok := m.buffer[tag]; ok {
		for _, msg := range record.Messages {
			err := protocol.receive(msg)
			if err != nil {
				jww.WARN.Printf("Failed to receive message: %+v", msg)
			}
		}
		delete(m.buffer, tag)
	}

	// create the long running thread which cleans the fingerprint list
	go func() {
		ticker := time.NewTicker(2 * flags.MaxGossipAge)
		for {
			<-ticker.C
			protocol.swapFingerprint()
		}
	}()

	m.bufferLock.Unlock()
}

// Returns the Gossip object for the provided tag from the Manager
func (m *Manager) Get(tag string) (*Protocol, bool) {
	m.protocolLock.RLock()
	defer m.protocolLock.RUnlock()

	p, ok := m.protocols[tag]
	return p, ok
}

// Deletes a Protocol from the Manager
func (m *Manager) Delete(tag string) {
	m.protocolLock.Lock()
	defer m.protocolLock.Unlock()

	delete(m.protocols, tag)
}

// Long-running thread to delete any messages in buffer older than 5m
func (m *Manager) bufferMonitor() chan bool {
	killChan := make(chan bool, 0)
	bufferExpirationTime := m.flags.BufferExpirationTime // Time in seconds that a record in the buffer should last
	frequency := m.flags.MonitorThreadFrequency

	go func() {
		for {
			// Loop through each tag in the buffer
			m.bufferLock.Lock()
			for tag, record := range m.buffer {
				if time.Since(record.Timestamp) > bufferExpirationTime {
					delete(m.buffer, tag)
				}
			}
			m.bufferLock.Unlock()

			select {
			case <-killChan:
				return
			default:
				time.Sleep(frequency)
			}
		}
	}()

	return killChan
}

const WorkerTimeout = 3 * time.Second

// launches numWorkers routines to handle sending of gossips for this protocol
// launches numWorkers routines to handle sending of gossips for this protocol
func launchSendWorkers(numWorkers uint32, receiver chan sendInstructions) {
	for i := uint32(0); i < numWorkers; i++ {
		go func() {
			for {
				// get a gossip send
				instructions := <-receiver
				// do the send
				errChan := make(chan error)

				go func() {
					err := instructions.sendFunc(instructions.peer)
					select {
					case errChan <- err:
					default:
						jww.WARN.Printf("Failed to send error report to "+
							"source for send to %s: %+v", instructions.peer, err)
					}
				}()

				var err error

				select {
				case err = <-errChan:
				case <-time.After(WorkerTimeout):
					jww.WARN.Printf("Send to peer %s failed after %s\n"+
						"\t ",
						instructions.peer, WorkerTimeout)
				}

				// handle errors if they occur
				if err != nil {
					select {
					case instructions.errChannel <- errors.WithMessagef(err,
						"Failed to send to ID %s", instructions.peer):
					default:
						jww.WARN.Print("Could not transmit gossip error")
					}
				}
				// signal the wait group
				instructions.wait.Done()
			}
		}()
	}
}
