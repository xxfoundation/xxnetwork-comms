////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/large"
	"gitlab.com/elixxir/primitives/ndf"
	"reflect"
	"testing"
)

// Happy path
func TestGroup_Get(t *testing.T) {
	// Create group with this group string
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	ourGroup := &Group{
		groupString: expectedGroup.String(),
	}

	// Fetch the group string
	receivedGroup := ourGroup.Get()

	// Compare the received group with the expected group
	if !reflect.DeepEqual(expectedGroup.String(), receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", expectedGroup.String(), receivedGroup)
	}
}

// Happy path
func TestNewGroup(t *testing.T) {
	// Create group
	ourNewGrp := NewGroup()
	// Check that the values are nil upon creation
	if ourNewGrp.groupString != "" || ourNewGrp.cyclicGroup != nil {
		t.Errorf("Values within group expected to be uninitialized upon creation."+
			"\n\tGroupString: %+v"+
			"\n\tCyclic Group: %+v", ourNewGrp.groupString, ourNewGrp.cyclicGroup)
	}

}

// Happy path
func TestGroup_Update(t *testing.T) {
	// Create group
	ourNewGrp := NewGroup()

	// Update the group with this group
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}
	err := ourNewGrp.Update(expectedGroup.String())
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// Check grpString
	expectedString := expectedGroup.String()
	if ourNewGrp.groupString != expectedString {
		t.Errorf("Update did not create expected string."+
			"\n\tExpected: %+v"+
			"\n\tCyclic Group: %+v", expectedString, ourNewGrp.cyclicGroup)
	}

	// Check cyclic.Group object
	expectedCyclic := cyclic.NewGroup(large.NewIntFromString("123", 16), large.NewIntFromString("2", 16))
	if !reflect.DeepEqual(expectedCyclic, ourNewGrp.cyclicGroup) {
		t.Errorf("Update did not produce expected cyclic.Group object."+
			"\n\tExpected: %+v"+
			"\n\tCyclic Group: %+v", expectedCyclic, ourNewGrp.cyclicGroup)

	}
}

// Happy path
func TestGroup_Update_DoubleUpdate(t *testing.T) {
	// Create group
	ourNewGrp := NewGroup()

	// Set up a group
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	// Update values to be initia
	err := ourNewGrp.Update(expectedGroup.String())
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// Attempt to update again with the same group string.
	// Should not error as we are not trying to change the initialized groups
	err = ourNewGrp.Update(expectedGroup.String())
	if err != nil {
		t.Errorf("Should not error when calling update with same value: %+v", err)
	}

	// Check grpString
	expectedString := expectedGroup.String()
	if ourNewGrp.groupString != expectedString {
		t.Errorf("Update did not create expected string."+
			"\n\tExpected: %+v"+
			"\n\tCyclic Group: %+v", expectedString, ourNewGrp.cyclicGroup)
	}

	// Check cyclic.Group object
	expectedCyclic := cyclic.NewGroup(large.NewIntFromString("123", 16), large.NewIntFromString("2", 16))
	if !reflect.DeepEqual(expectedCyclic, ourNewGrp.cyclicGroup) {
		t.Errorf("Update did not produce expected cyclic.Group object."+
			"\n\tExpected: %+v"+
			"\n\tCyclic Group: %+v", expectedCyclic, ourNewGrp.cyclicGroup)

	}

}

// Error path
func TestGroup_Update_DoubleUpdate_Error(t *testing.T) {
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	ourNewGrp := NewGroup()

	err := ourNewGrp.Update(expectedGroup.String())
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// A group that does not match the initialized group
	badGroup := ndf.Group{
		Prime:      "69",
		SmallPrime: "420",
		Generator:  "98",
	}

	// Attempt to update again with a different group string.
	// Should error as you should not be able to update values once initialized
	err = ourNewGrp.Update(badGroup.String())
	if err != nil {
		return
	}

	t.Errorf("Expected error case: Should error when trying to modify values in group!")

}
