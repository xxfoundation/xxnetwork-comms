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
	errorCounter *uint64
}

// Constructor for a Metric object
func NewMetric() *Metric {
	errCounter := uint64(0)
	return &Metric{
		errorCounter: &errCounter,
	}
}

// Returns a copy of Metric and resets internal state
func (m *Metric) Get() *Metric {
	metricCopy := m.deepCopy()
	atomic.StoreUint64(m.errorCounter, 0)
	return metricCopy
}

// Increments the error counter in a thread-safe manner
func (m *Metric) IncrementErrors() {
	atomic.AddUint64(m.errorCounter, 1)
}

// DeepCopy creates a copy of Metric.
func (m *Metric) deepCopy() *Metric {
	newErrCounter := atomic.LoadUint64(m.errorCounter)
	return &Metric{
		errorCounter: &newErrCounter,
	}
}
