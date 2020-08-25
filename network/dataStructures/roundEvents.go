///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Stores callbacks that will be called in the process of running a round
package dataStructures

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/primitives/id"
	"sync"
	"time"
)

// Callbacks must use this function signature
type RoundEventCallback func(ri *pb.RoundInfo, timedOut bool)

// One callback and associated data
type EventCallback struct {
	// Round states where this function can be called
	states []states.Round
	// Send on this channel to cause the relevant callbacks
	signal chan *pb.RoundInfo
}

// Holds the callbacks for a round
type RoundEvents struct {
	// The slice that map[id.Round] maps to is a collection of event callbacks for each of the round's states
	callbacks map[id.Round][states.NUM_STATES]map[*EventCallback]*EventCallback
	mux       sync.RWMutex
}

func NewRoundEvents() *RoundEvents {
	return &RoundEvents{
		callbacks: make(map[id.Round][states.NUM_STATES]map[*EventCallback]*EventCallback),
	}
}

func (r *RoundEvents) Remove(rid id.Round, e *EventCallback) {
	r.mux.Lock()
	r.remove(rid, e)
	r.mux.Unlock()
}

// Remove an event callback from all the states' maps
func (r *RoundEvents) remove(rid id.Round, e *EventCallback) {
	for _, s := range e.states {
		delete(r.callbacks[rid][s], e)
	}

	// Remove this round's events from the top-level map if there aren't any
	// callbacks left in any of the states
	removeRound := true
	for s := states.Round(0); (s < states.NUM_STATES) && removeRound; s++ {
		removeRound = removeRound && len(r.callbacks[rid][s]) == 0
	}
	if removeRound {
		delete(r.callbacks, rid)
	}
}

func (r *RoundEvents) signal(rid id.Round, event *EventCallback, callback RoundEventCallback, timeout time.Duration) {
	ri := &pb.RoundInfo{ID: uint64(rid)}
	select {
	case <-time.After(timeout):
		go r.Remove(rid, event)
		callback(ri, true)
	case ri = <-event.signal:
		callback(ri, false)
	}
}

// Adds an event to the RoundEvents struct and returns its handle for possible deletion
func (r *RoundEvents) AddRoundEvent(rid id.Round, callback RoundEventCallback, timeout time.Duration, validStates ...states.Round) *EventCallback {
	// Add the specific event to the round
	thisEvent := &EventCallback{
		states: validStates,
		signal: make(chan *pb.RoundInfo, 1),
	}

	go r.signal(rid, thisEvent, callback, timeout)

	r.mux.Lock()
	callbacks, ok := r.callbacks[rid]
	if !ok {
		// create callbacks for this round
		for i := range callbacks {
			callbacks[i] = make(map[*EventCallback]*EventCallback)
		}

		r.callbacks[rid] = callbacks
	}

	for _, s := range validStates {
		callbacks[s][thisEvent] = thisEvent
	}
	r.mux.Unlock()
	return thisEvent
}

func (r *RoundEvents) TriggerRoundEvent(ri *pb.RoundInfo) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	// Try to find callbacks
	callbacks, ok := r.callbacks[id.Round(ri.ID)]
	if !ok {
		return
	}
	thisStatesCallbacks := callbacks[ri.State]
	if len(thisStatesCallbacks) != 0 {
		// Keep track of events we've used for later removal
		var events []*EventCallback
		for _, event := range thisStatesCallbacks {
			event.signal <- ri
			events = append(events, event)
		}
		// Everything we sent a signal to is no longer needed
		for _, event := range events {
			r.remove(id.Round(ri.ID), event)
		}
	}
}
