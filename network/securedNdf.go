////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper for the ndf object

package network

import (
	"fmt"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/ndf"
)

// wraps the ndf data structure, expoting all the functionality expect the
// ability to change the ndf
type SecuredNdf struct {
	f *ds.Ndf
}

// Initialize a securedNdf from a primitives NetworkDefinition object
func NewSecuredNdf(definition *ndf.NetworkDefinition) (*SecuredNdf, error) {
	ndf, err := ds.NewNdf(definition)
	if err != nil {
		return nil, err
	}
	return &SecuredNdf{
		ndf,
	}, nil
}

// unexported NDF update code
func (sndf *SecuredNdf) update(m *pb.NDF, key *rsa.PublicKey) error {
	err := signature.Verify(m, key)
	if err != nil {
		fmt.Printf("err: %+v", err)
		return errors.WithMessage(err, "Could not validate NDF")
	}

	return sndf.f.Update(m)
}

// Get the primitives object for an ndf
func (sndf *SecuredNdf) Get() *ndf.NetworkDefinition {
	return sndf.f.Get()
}

// get the hash of the ndf
func (sndf *SecuredNdf) GetHash() []byte {
	return sndf.f.GetHash()
}

// get the protobuf message NDF
func (sndf *SecuredNdf) GetPb() *pb.NDF {
	return sndf.f.GetPb()
}

// Compare a hash to the stored
func (sndf *SecuredNdf) CompareHash(h []byte) bool {
	return sndf.f.CompareHash(h)
}
