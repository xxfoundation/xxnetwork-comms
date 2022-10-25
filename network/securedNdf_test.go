////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package network

import (
	"fmt"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/comms/testutils"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/ndf"
	"math/rand"
	"testing"
)

func setup() *ds.Ndf {
	msg := &pb.NDF{
		Ndf: testutils.ExampleNDF,
	}
	netDef := &ds.Ndf{}

	_ = netDef.Update(msg)
	return netDef
}

func TestNewSecuredNdf(t *testing.T) {
	d, _ := ndf.Unmarshal(testutils.ExampleNDF)
	sndf, err := NewSecuredNdf(d)
	if err != nil {
		t.Error(err)
	}
	if sndf == nil {
		t.Errorf("Internal ndf object is nil")
	}
}

func TestSecuredNdf_update(t *testing.T) {
	src := rand.New(rand.NewSource(42))
	privKey, err := rsa.GenerateKey(src, 768)
	if err != nil {
		t.Errorf("Could not generate rsa key: %s", err)
	}

	badSrc := rand.New(rand.NewSource(33))
	badPriv, err := rsa.GenerateKey(badSrc, 768)
	if err != nil {
		t.Errorf("Could not generate rsa key: %s", err)
	}
	badPub := badPriv.GetPublic()
	fmt.Println(badPub)

	f := pb.NDF{}

	baseNDF := ndf.NetworkDefinition{}
	f.Ndf, err = baseNDF.Marshal()

	if err != nil {
		t.Errorf("Could not generate serialized ndf: %s", err)
	}

	err = signature.SignRsa(&f, privKey)

	if err != nil {
		t.Errorf("Could not sign serialized ndf: %s", err)
	}

	sndf, err := NewSecuredNdf(testutils.NDF)
	if err != nil {
		t.Errorf("Failed to secure ndf: %+v", err)
	}
	err = sndf.update(&f, privKey.GetPublic())

	if err != nil {
		t.Errorf("Could not update ndf: %s", err)
	}

	err = sndf.update(&f, badPub)
	// Fixme
	/*	if err == nil {
		t.Errorf("should have received bad key error")
	}*/

}

func TestSecuredNdf_Get(t *testing.T) {
	sn := SecuredNdf{f: setup()}
	if sn.Get() == nil {
		t.Error("Should have received ndf")
	}
}

func TestSecuredNdf_GetPb(t *testing.T) {
	sn := SecuredNdf{f: setup()}
	if sn.GetPb() == nil {
		t.Error("Should have received pb")
	}
}

func TestSecuredNdf_GetHash(t *testing.T) {
	sn := SecuredNdf{f: setup()}
	if sn.GetHash() == nil {
		t.Error("Should have received hash")
	}
}

func TestSecuredNdf_CompareHash(t *testing.T) {
	sn := SecuredNdf{f: setup()}
	b := sn.CompareHash(sn.f.GetHash())
	if !b {
		t.Error("Should have received true for comparison")
	}
}
