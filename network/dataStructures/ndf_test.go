////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testutils"
	"testing"
)

func setup() *Ndf {
	msg := &mixmessages.NDF{
		Ndf: testutils.ExampleNDF,
	}
	ndf := &Ndf{}

	_ = ndf.Update(msg)
	return ndf
}

func TestNdf_Get(t *testing.T) {
	ndf := setup()

	if ndf.Get() == nil {
		t.Error("Should have returned ndf.f")
	}
}

func TestNdf_Update(t *testing.T) {
	msg := &mixmessages.NDF{
		Ndf: testutils.ExampleNDF,
	}
	badMsg := &mixmessages.NDF{
		Ndf: []byte("lasagna"),
	}
	ndf := Ndf{}

	err := ndf.Update(badMsg)
	if err == nil {
		t.Error("Should have returned error when unable to decode ndf")
	}

	err = ndf.Update(msg)
	if err != nil {
		t.Errorf("Failed to update ndf: %+v", err)
	}

	if ndf.f == nil || ndf.hash == nil || ndf.pb == nil {
		t.Error("Failed to properly set contents of ndf object")
	}
}

func TestNdf_GetHash(t *testing.T) {
	ndf := setup()

	if ndf.GetHash() == nil {
		t.Error("Hash not properly returned")
	}
}

func TestNdf_GetPb(t *testing.T) {
	ndf := setup()

	if ndf.GetPb() == nil {
		t.Error("Pb not properly set")
	}
}

func TestNdf_CompareHash(t *testing.T) {
	ndf := &Ndf{}

	ndf = setup()
	b := ndf.CompareHash(ndf.hash)
	if !b {
		t.Error("Should return true when hashes are the same")
	}

	b = ndf.CompareHash([]byte("test"))
	if b {
		t.Error("Should return false when hashes are different")
	}
}
