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
	metric := NewMetric()

	expectedErrCnt := uint64(0)
	expectedMetric := &Metric{
		errorCounter: &expectedErrCnt,
	}

	if !reflect.DeepEqual(expectedMetric, metric) {
		t.Errorf("Unexpected values in constructed Metric object."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedMetric, metric)
	}
}

// Unit test of deepCopy
func TestMetric_deepCopy(t *testing.T) {
	metric := NewMetric()

	metricCopy := metric.deepCopy()

	if !reflect.DeepEqual(metricCopy, metric) {
		t.Errorf("Deep copy did not create identical copy."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", metric, metricCopy)
	}
}

// Unit test for IncrementErrors
func TestMetric_IncrementErrors(t *testing.T) {
	expectedCount := 25
	metric := NewMetric()

	for i := 0; i < expectedCount; i++ {
		metric.IncrementErrors()
	}

	if *metric.errorCounter != uint64(expectedCount) {
		t.Errorf("IncrementErrors did not "+
			"result in expected error count."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedCount, *metric.errorCounter)
	}
}

// Unit test for Get()
func TestMetric_Get(t *testing.T) {
	expectedCount := 25
	metric := NewMetric()

	for i := 0; i < expectedCount; i++ {
		metric.IncrementErrors()
	}

	// Check that the metricCopy has the expected error count
	metricCopy := metric.Get()
	if *metricCopy.errorCounter != uint64(expectedCount) {
		t.Errorf("Get() did not pull expected state."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedCount, *metricCopy.errorCounter)
	}

	// Check that the original metric's state has been reset
	if *metric.errorCounter != uint64(0) {
		t.Errorf("Get call should reset state for metric")
	}
}
