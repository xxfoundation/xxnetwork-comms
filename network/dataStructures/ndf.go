///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Handles basic operations on different forms of network definitions

package dataStructures

import (
	"bytes"
	"crypto/sha256"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/primitives/ndf"
	"sync"
)

// Struct which encapsulates all data from an NDF
type Ndf struct {
	f    *ndf.NetworkDefinition
	pb   *pb.NDF
	hash []byte
	sync.RWMutex
}

// Initialize an Ndf object from a primitives NetworkDefinition
func NewNdf(definition *ndf.NetworkDefinition) (*Ndf, error) {
	h, err := GenerateNDFHash(definition)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to hash ndf")
	}
	return &Ndf{
		f:    definition,
		pb:   nil,
		hash: h,
	}, nil
}

//Updates to a new NDF if the passed NDF is valid
func (file *Ndf) Update(m *pb.NDF) error {

	//build the ndf object
	decoded, _, err := ndf.DecodeNDF(string(m.Ndf))

	if err != nil {
		return errors.WithMessage(err, "Could not decode the NDF")
	}

	file.Lock()
	defer file.Unlock()

	file.pb = m
	file.f = decoded

	file.hash, err = GenerateNDFHash(file.f)

	return err
}

//returns the ndf object
//fix-me: return a copy instead to ensure edits to not impact the
//original version
func (file *Ndf) Get() *ndf.NetworkDefinition {
	file.RLock()
	defer file.RUnlock()

	return file.f
}

//returns the ndf hash
func (file *Ndf) GetHash() []byte {
	file.RLock()
	defer file.RUnlock()

	rtn := make([]byte, len(file.hash))
	copy(rtn, file.hash)
	return rtn
}

//returns the ndf hash
//fix-me: return a copy instead to ensure edits to not impact the
//original version
func (file *Ndf) GetPb() *pb.NDF {
	file.RLock()
	defer file.RUnlock()

	return file.pb
}

// Evaluates if the passed ndf hash is the same as the stored one
func (file *Ndf) CompareHash(h []byte) bool {
	file.RLock()
	defer file.RUnlock()

	// Return whether the hashes are different
	return bytes.Compare(file.hash, h) == 0
}

// helper function to generate a hash of the NDF
func GenerateNDFHash(definition *ndf.NetworkDefinition) ([]byte, error) {
	//set the ndf hash
	marshaled, err := definition.Marshal()
	if err != nil {
		return nil, errors.WithMessage(err,
			"Could not marshal NDF for hashing")
	}

	// Serialize then hash the constructed ndf
	hash := sha256.New()
	hash.Write(marshaled)
	return hash.Sum(nil), nil
}
