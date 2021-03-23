///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"reflect"
	"testing"
)

// Smoke test for constructor
func TestNewMetric(t *testing.T) {
	metric := newMetric()

	expectedErrCnt := uint64(0)
	expectedMetric := &Metric{
		ErrCounter: &expectedErrCnt,
	}

	if !reflect.DeepEqual(expectedMetric, metric) {
		t.Errorf("Unexpected values in constructed Metric object."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedMetric, metric)
	}
}

// Unit test of deepCopy
func TestMetric_deepCopy(t *testing.T) {
	metric := newMetric()

	metricCopy := metric.deepCopy()

	if !reflect.DeepEqual(metricCopy, metric) {
		t.Errorf("Deep copy did not create identical copy."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", metric, metricCopy)
	}
}

// Unit test for incrementErrors
func TestMetric_IncrementErrors(t *testing.T) {
	expectedCount := 25
	metric := newMetric()

	for i := 0; i < expectedCount; i++ {
		metric.incrementErrors()
	}

	if *metric.ErrCounter != uint64(expectedCount) {
		t.Errorf("incrementErrors did not "+
			"result in expected error count."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedCount, *metric.ErrCounter)
	}
}

// Unit test for get()
func TestMetric_Get(t *testing.T) {
	expectedCount := 25
	metric := newMetric()

	for i := 0; i < expectedCount; i++ {
		metric.incrementErrors()
	}

	// Check that the metricCopy has the expected error count
	metricCopy := metric.get()
	if *metricCopy.ErrCounter != uint64(expectedCount) {
		t.Errorf("get() did not pull expected state."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedCount, *metricCopy.ErrCounter)
	}

	// Check that the original metric's state has been reset
	if *metric.ErrCounter != uint64(0) {
		t.Errorf("get call should reset state for metric")
	}
}
