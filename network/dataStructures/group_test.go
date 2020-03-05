////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"gitlab.com/elixxir/primitives/ndf"
	"reflect"
	"testing"
)

// Happy path
func TestGroup_Get(t *testing.T) {
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	ourGroup := &Group{
		grp: expectedGroup.String(),
	}

	receivedGroup := ourGroup.Get()

	if !reflect.DeepEqual(expectedGroup.String(), receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup.String(), receivedGroup)
	}
}

// Happy path
func TestNewGroup(t *testing.T) {
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	recievedGroup := NewGroup(expectedGroup.String()).Get()

	if !reflect.DeepEqual(expectedGroup.String(), recievedGroup) {
		t.Errorf("\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup.String(), recievedGroup)
	}
}
