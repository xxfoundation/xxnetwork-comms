package network

import "testing"

func TestNewSecuredNdf(t *testing.T) {
	sndf := NewSecuredNdf()
	if sndf==nil{
		t.Errorf("Internal ndf object is nil")
	}
}

func TestSecuredNdf_CompareHash(t *testing.T) {

}