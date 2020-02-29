package network

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/ndf"
	"math/rand"
	"testing"
)

func TestNewSecuredNdf(t *testing.T) {
	sndf := NewSecuredNdf()
	if sndf==nil{
		t.Errorf("Internal ndf object is nil")
	}
}

func TestSecuredNdf_update(t *testing.T) {
	src:=rand.New(rand.NewSource(42))
	privKey, err := rsa.GenerateKey(src, 768)

	if err!=nil{
		t.Errorf("Could not generate rsa key: %s", err)
	}

	f := pb.NDF{}

	baseNDF := ndf.NetworkDefinition{}
	f.Ndf, err = baseNDF.Marshal()

	if err!=nil{
		t.Errorf("Could not generate serialized ndf: %s", err)
	}

	err = signature.Sign(&f, privKey)

	if err!=nil{
		t.Errorf("Could not sign serialized ndf: %s", err)
	}

	sndf := NewSecuredNdf()
	err = sndf.update(&f,privKey.GetPublic())

	if err!=nil{
		t.Errorf("Could not update ndf: %s", err)
	}

}