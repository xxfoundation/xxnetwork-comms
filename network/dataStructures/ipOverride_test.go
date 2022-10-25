////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

//tests that NewIpOverrideList returns a properly formatted override list
func TestNewIpOverrideList(t *testing.T) {
	nol := NewIpOverrideList()

	if nol.ipOverride == nil {
		t.Errorf("Ip Override List should have a map")
	}
}

// tests that creating an override works
func TestIpOverrideList_Override(t *testing.T) {
	iol := &IpOverrideList{
		ipOverride: make(map[id.ID]string),
	}

	testID := id.NewIdFromUInt(42, id.Node, t)
	testIP := "woop"

	iol.Override(testID, testIP)

	if len(iol.ipOverride) != 1 {
		t.Errorf("IP override has the wrong length")
	}

	resultIP, exist := iol.ipOverride[*testID]

	if !exist {
		t.Errorf("could not find override in the map")
	}

	if resultIP != testIP {
		t.Errorf("the ip returned from the map is not as expected"+
			" expected: %s, recieved: %s", testIP, resultIP)
	}
}

// tests that the old IP is passed through when there is no override
func TestIpOverrideList_CheckOverride_NoOverride(t *testing.T) {
	iol := &IpOverrideList{
		ipOverride: make(map[id.ID]string),
	}

	testID := id.NewIdFromUInt(42, id.Node, t)
	testIP := "woop"

	dummyID := id.NewIdFromUInt(69, id.Node, t)
	iol.Override(dummyID, "blarg")

	resultIP := iol.CheckOverride(testID, testIP)
	if resultIP != testIP {
		t.Errorf("The returned ip is not as expected; "+
			"Expected: %s, Returned: %s", testIP, resultIP)
	}
}

// tests that the old IP overwritten through when there is an override
func TestIpOverrideList_CheckOverride_ValidOverride(t *testing.T) {
	iol := &IpOverrideList{
		ipOverride: make(map[id.ID]string),
	}

	testID := id.NewIdFromUInt(42, id.Node, t)
	testIP := "woop"

	iol.Override(testID, testIP)

	resultIP := iol.CheckOverride(testID, "blarg")
	if resultIP != testIP {
		t.Errorf("The returned ip is not as expected; "+
			"Expected: %s, Returned: %s", testIP, resultIP)
	}
}
