///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functionality for metric tracking on sending

package connect

import "sync/atomic"

type Metric struct {
	// Active count of non-excluded errors
	// (ie errors we deem important)
	ErrCounter *uint64
}

// Constructor for a Metric object
func newMetric() *Metric {
	errCounter := uint64(0)
	return &Metric{
		ErrCounter: &errCounter,
	}
}

// Returns a copy of Metric and resets internal state
func (m *Metric) get() *Metric {
	metricCopy := m.deepCopy()
	atomic.StoreUint64(m.ErrCounter, 0)
	return metricCopy
}

// Increments the error counter in a thread-safe manner
func (m *Metric) incrementErrors() {
	atomic.AddUint64(m.ErrCounter, 1)
}

// deepCopy creates a copy of Metric.
func (m *Metric) deepCopy() *Metric {
	newErrCounter := atomic.LoadUint64(m.ErrCounter)
	return &Metric{
		ErrCounter: &newErrCounter,
	}
}
