////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/ndf"
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

	gstr, err := expectedGroup.String()

	if err != nil {
		t.Errorf("Could not get expected group's string: %s", err)
	}
	ourGroup := &Group{
		groupString: gstr,
	}

	// Fetch the group string
	receivedGroup := ourGroup.GetString()

	// Compare the received group with the expected group
	if !reflect.DeepEqual(gstr, receivedGroup) {
		t.Errorf("Getter didn't get expected value! "+
			"\n\tExpected: %+v"+
			"\n\tReceived: %+v", gstr, receivedGroup)
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

	// Set up a ndf.group
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	gstr, err := expectedGroup.String()

	if err != nil {
		t.Errorf("Could not get expected group's string: %s", err)
	}

	// Update values to be initialized
	err = ourNewGrp.Update(gstr)
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// Check grpString
	if ourNewGrp.groupString != gstr {
		t.Errorf("Update did not create expected string."+
			"\n\tExpected: %+v"+
			"\n\tCyclic Group: %+v", gstr, ourNewGrp.cyclicGroup)
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

	// Set up a ndf.group
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	gstr, err := expectedGroup.String()

	if err != nil {
		t.Errorf("Could not get expected group's string: %s", err)
	}

	// Update values to be initialized
	err = ourNewGrp.Update(gstr)
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// Attempt to update again with the same group string.
	// Should not error as we are not trying to change the initialized groups
	err = ourNewGrp.Update(gstr)
	if err != nil {
		t.Errorf("Should not error when calling update with same value: %+v", err)
	}

	// Check grpString
	expectedString := gstr
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
	// Create group
	ourNewGrp := NewGroup()

	// Set up a group
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	gstr, err := expectedGroup.String()

	if err != nil {
		t.Errorf("Could not get expected group's string: %s", err)
	}

	// Update values to be initialized
	err = ourNewGrp.Update(gstr)
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// A group that does not match the initialized group
	badGroup := ndf.Group{
		Prime:      "69",
		SmallPrime: "420",
		Generator:  "98",
	}

	bgstr, err := badGroup.String()

	if err != nil {
		t.Errorf("Could not get expected group's string: %s", err)
	}

	// Attempt to update again with a different group string.
	// Should error as you should not be able to update values once initialized
	err = ourNewGrp.Update(bgstr)
	if err != nil {
		return
	}

	t.Errorf("Expected error case: Should error when trying to modify values in group!")

}
