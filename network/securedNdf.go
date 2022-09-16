////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Wrapper for the ndf object

package network

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/ndf"
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
	err := signature.VerifyRsa(m, key)
	if err != nil {
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
