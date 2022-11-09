////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package dataStructures stores callbacks that will be called in the process of
// running a round.
package dataStructures

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/primitives/id"
	"sync"
	"time"
)

// RoundEventCallback is the callbacks called on trigger.
type RoundEventCallback func(ri *pb.RoundInfo, timedOut bool)

// EventCallback contains one callback and associated data.
type EventCallback struct {
	// Round states where this function can be called
	states []states.Round

	// Send on this channel to cause the relevant callbacks
	signal chan *pb.RoundInfo
}

// RoundEvents holds the callbacks for a round.
type RoundEvents struct {
	// The slice that map[id.Round] maps to is a collection of event callbacks
	// for each of the round's states
	callbacks map[id.Round][states.NUM_STATES]map[*EventCallback]*EventCallback
	mux       sync.RWMutex
}

// NewRoundEvents initialize a new RoundEvents object.
func NewRoundEvents() *RoundEvents {
	return &RoundEvents{
		callbacks: make(
			map[id.Round][states.NUM_STATES]map[*EventCallback]*EventCallback),
	}
}

// Remove wraps non-exported remove with mutex.
func (r *RoundEvents) Remove(rid id.Round, e *EventCallback) {
	r.mux.Lock()
	r.remove(rid, e)
	r.mux.Unlock()
}

// remove deletes an event callback from all the states' maps. Also removes the
// round from the top-level map if it becomes empty.
func (r *RoundEvents) remove(rid id.Round, e *EventCallback) {
	// Remove all EventCallbacks for the given round
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

// signal calls or timeout a round event. Removes round events when they are
// called or timed out to allow them to get garbage collected.
func (r *RoundEvents) signal(rid id.Round, event *EventCallback,
	callback RoundEventCallback, timeout time.Duration) {
	ri := &pb.RoundInfo{ID: uint64(rid)}
	select {
	case <-time.After(timeout):
		go r.Remove(rid, event)
		callback(ri, true)
	case ri = <-event.signal:
		go r.Remove(rid, event)
		callback(ri, false)
	}
}

type EventReturn struct {
	RoundInfo *pb.RoundInfo
	TimedOut  bool
}

// AddRoundEventChan puts the round event on a channel instead of using a
// callback.
func (r *RoundEvents) AddRoundEventChan(rid id.Round,
	eventChan chan EventReturn, timeout time.Duration,
	validStates ...states.Round) *EventCallback {

	callback := func(ri *pb.RoundInfo, timedOut bool) {
		eventChan <- EventReturn{ri, timedOut}
	}

	return r.AddRoundEvent(rid, callback, timeout, validStates...)
}

// AddRoundEvent adds an event to the RoundEvents struct and returns its handle
// for possible deletion.
func (r *RoundEvents) AddRoundEvent(rid id.Round, callback RoundEventCallback,
	timeout time.Duration, validStates ...states.Round) *EventCallback {
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

// TriggerRoundEvent signals all round events matching the passed RoundInfo
// according to its ID and state.
func (r *RoundEvents) TriggerRoundEvent(rnd *Round) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	// Try to find callbacks
	callbacks, ok := r.callbacks[id.Round(rnd.info.ID)]
	if !ok || len(callbacks[rnd.info.State]) == 0 {
		return
	}

	// Retrieve and validate the round info
	roundInfo := rnd.Get()

	// Send round info to every event in the list
	for _, event := range callbacks[rnd.info.State] {
		select {
		case event.signal <- roundInfo:
		default:
		}
	}
}

// TriggerRoundEvents signals all round events matching the passed RoundInfos
// according to its ID and state.
func (r *RoundEvents) TriggerRoundEvents(rounds ...*Round) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	for _, rnd := range rounds {

		// Try to find callbacks
		callbacks, ok := r.callbacks[id.Round(rnd.info.ID)]
		if !ok || len(callbacks[rnd.info.State]) == 0 {
			continue
		}

		// Retrieve and validate the round info
		roundInfo := rnd.Get()

		// Send round info to every event in the list
		for _, event := range callbacks[rnd.info.State] {
			select {
			case event.signal <- roundInfo:
			default:
			}
		}
	}
}
