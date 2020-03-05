////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/large"
	"reflect"
	"sync"
	"testing"
)

// Happy path
func TestGroup_Get(t *testing.T) {
	p := large.NewInt(33)
	g := large.NewInt(29)
	expectedGroup := cyclic.NewGroup(p, g)
	ourGroup := &Group{
		grp:   expectedGroup,
		mutex: &sync.RWMutex{},
	}

	receivedGroup := ourGroup.Get()

	if !reflect.DeepEqual(expectedGroup, receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup, receivedGroup)
	}
}

// Happy path
func TestNewGroup(t *testing.T) {
	p := large.NewInt(33)
	g := large.NewInt(29)

	expectedGroup := cyclic.NewGroup(p, g)

	recievedGroup := NewGroup(p, g).Get()

	if !reflect.DeepEqual(expectedGroup, recievedGroup) {
		t.Errorf("\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup, recievedGroup)
	}
}
