///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Stores callbacks that will be called in the process of running a round
package dataStructures

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/states"
	"sync"
	"time"
)

// One callback and associated data
type eventCallback struct {
	// Function that will be called
	f RoundEventCallback
	// Round states where this function can be called
	states []states.Round
	// Callback won't be called after this time
	timeout time.Time
}

// Has all the callbacks and information for a round
type roundEventCallbacks struct {
	// All callbacks associated with the round
	callbacks []eventCallback
	// Time at which all callbacks will have timed out and this round will be removed from the map
	deletionTime time.Time
}

// Callbacks must use this function signature
type RoundEventCallback func(ri *pb.RoundInfo)

// Holds the callbacks for a round
type RoundEvents struct {
	// Callbacks
	callbacks map[id.Round]*roundEventCallbacks
	mux       sync.RWMutex
}

// Initialize the round events structure
func NewRoundEvents() *RoundEvents {
	return &RoundEvents{
		callbacks: make(map[id.Round]*roundEventCallbacks),
	}
}

// Add a round event callback for a certain round and states
// If the callback is called after a certain timeout
func (r *RoundEvents) AddRoundEvent(id *id.Round, callback RoundEventCallback, timeout time.Duration, states ...states.Round) {
	r.mux.Lock()
	callbacks, ok := r.callbacks[*id]
	if !ok {
		// create callbacks for this round
		callbacks = &roundEventCallbacks{
			deletionTime: time.Now(),
		}
		r.callbacks[*id] = callbacks
	}

	// Add the specific event to the round
	thisEvent := eventCallback{
		f:       callback,
		states:  states,
		timeout: time.Now().Add(timeout),
	}
	callbacks.callbacks = append(callbacks.callbacks, thisEvent)

	// Reschedule deletion if this event's timeout is after the callback structure's deletion time
	if callbacks.deletionTime.Before(thisEvent.timeout) {
		callbacks.deletionTime = thisEvent.timeout
		time.AfterFunc(timeout, func() {
			now := time.Now()
			r.mux.Lock()
			// Check to see whether it's time to remove this round or not
			// which is why it's necessary to check)
			// (another call to AddRoundEvent could have changed the deletion time,
			if now.After(callbacks.deletionTime) || now.Equal(callbacks.deletionTime) {
				// OK to delete this round's callbacks from the map
				delete(r.callbacks, *id)
			}
			r.mux.Unlock()
		})
	}
	r.mux.Unlock()
}

// Returns error if no event is found
func (r *RoundEvents) TriggerRoundEvent(ri *pb.RoundInfo) error {
	r.mux.RLock()
	// Try to find callbacks
	callbacks, ok := r.callbacks[id.Round(ri.ID)]
	if !ok {
		r.mux.RUnlock()
		return errors.New("no callbacks found for that round ID")
	}
	var wg sync.WaitGroup
	now := time.Now()
	foundCallback := false
	for i := range callbacks.callbacks {
		// launch all callbacks for this round's state
		for j := range callbacks.callbacks[i].states {
			// Callback must be relevant for this round state to be called
			if callbacks.callbacks[i].states[j] == states.Round(ri.State) {
				// Callback must not have timed out to be called
				if callbacks.callbacks[i].timeout.After(now) {
					wg.Add(1)
					foundCallback = true
					go func() {
						callbacks.callbacks[i].f(ri)
						wg.Done()
					}()
				}
			}
		}
	}
	wg.Wait()
	r.mux.RUnlock()

	if foundCallback {
		return nil
	} else {
		return errors.New("no callbacks found in that round's callbacks")
	}
}

// Remove all the events for a round so the memory can be garbage collected
// No effect if the specified round doesn't have events registered
func (r *RoundEvents) RemoveRoundEvents(id *id.Round) {
	r.mux.Lock()
	delete(r.callbacks, *id)
	r.mux.Unlock()
}
