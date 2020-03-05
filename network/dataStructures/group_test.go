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

/*
cmixGrp := cyclic.NewGroup(
large.NewIntFromString("9DB6FB5951B66BB6FE1E140F1D2CE5502374161FD6538DF1648218642F0B5C48"+
"C8F7A41AADFA187324B87674FA1822B00F1ECF8136943D7C55757264E5A1A44F"+
"FE012E9936E00C1D3E9310B01C7D179805D3058B2A9F4BB6F9716BFE6117C6B5"+
"B3CC4D9BE341104AD4A80AD6C94E005F4B993E14F091EB51743BF33050C38DE2"+
"35567E1B34C3D6A5C0CEAA1A0F368213C3D19843D0B4B09DCB9FC72D39C8DE41"+
"F1BF14D4BB4563CA28371621CAD3324B6A2D392145BEBFAC748805236F5CA2FE"+
"92B871CD8F9C36D3292B5509CA8CAA77A2ADFC7BFD77DDA6F71125A7456FEA15"+
"3E433256A2261C6A06ED3693797E7995FAD5AABBCFBE3EDA2741E375404AE25B", 16),
large.NewIntFromString("5C7FF6B06F8F143FE8288433493E4769C4D988ACE5BE25A0E24809670716C613"+
"D7B0CEE6932F8FAA7C44D2CB24523DA53FBE4F6EC3595892D1AA58C4328A06C4"+
"6A15662E7EAA703A1DECF8BBB2D05DBE2EB956C142A338661D10461C0D135472"+
"085057F3494309FFA73C611F78B32ADBB5740C361C9F35BE90997DB2014E2EF5"+
"AA61782F52ABEB8BD6432C4DD097BC5423B285DAFB60DC364E8161F4A2A35ACA"+
"3A10B1C4D203CC76A470A33AFDCBDD92959859ABD8B56E1725252D78EAC66E71"+
"BA9AE3F1DD2487199874393CD4D832186800654760E1E34C09E4D155179F9EC0"+
"DC4473F996BDCE6EED1CABED8B6F116F7AD9CF505DF0F998E34AB27514B0FFE7", 16))



*/
// Happy path
func TestGroup_Get(t *testing.T) {
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}

	ourGroup := &Group{
		groupString: expectedGroup.String(),
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

	ourNewGrp := NewGroup()

	if ourNewGrp.groupString != "" || ourNewGrp.cyclicGroup != nil {
		t.Errorf("Values within group expected to be uninitialized upon creation."+
			"\n\tGroupString: %+v"+
			"\n\tCyclic Group: %+v", ourNewGrp.groupString, ourNewGrp.cyclicGroup)
	}

}

// Happy path
func TestGroup_Update(t *testing.T) {
	expectedGroup := ndf.Group{
		Prime:      "123",
		SmallPrime: "456",
		Generator:  "2",
	}
	expectedString := expectedGroup.String()

	ourNewGrp := NewGroup()

	err := ourNewGrp.Update(expectedGroup.String())
	if err != nil {
		t.Errorf("Unable to update group: %+v", err)
	}

	// Check grpString
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

	//
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
