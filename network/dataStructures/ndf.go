////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles basic operations on different forms of network definitions

package dataStructures

import (
	"bytes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/primitives/ndf"
	"golang.org/x/crypto/blake2b"
	"sync"
)

// Ndf encapsulates all data from an NDF.
type Ndf struct {
	f    *ndf.NetworkDefinition
	pb   *pb.NDF
	hash []byte
	sync.RWMutex
}

// NewNdf initializes a Ndf object from a primitives ndf.NetworkDefinition.
func NewNdf(definition *ndf.NetworkDefinition) (*Ndf, error) {
	h, err := GenerateNDFHash(nil)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to hash NDF")
	}

	return &Ndf{
		f:    definition,
		pb:   nil,
		hash: h,
	}, nil
}

// Update to a new NDF if the passed NDF is valid.
func (file *Ndf) Update(m *pb.NDF) error {

	// Build the ndf object
	decoded, err := ndf.Unmarshal(m.Ndf)
	if err != nil {
		return errors.WithMessage(err, "Could not decode the NDF")
	}

	file.Lock()
	defer file.Unlock()

	file.pb = m
	file.f = decoded

	file.hash, err = GenerateNDFHash(file.pb)

	return err
}

// Get returns the NDF object.
// FIXME: return a copy instead to ensure edits to not impact the original version
func (file *Ndf) Get() *ndf.NetworkDefinition {
	file.RLock()
	defer file.RUnlock()

	return file.f
}

// GetHash returns the NDF hash.
func (file *Ndf) GetHash() []byte {
	file.RLock()
	defer file.RUnlock()

	rtn := make([]byte, len(file.hash))
	copy(rtn, file.hash)
	return rtn
}

// GetPb returns the NDF message.
// FIXME: return a copy instead to ensure edits to not impact the original version
func (file *Ndf) GetPb() *pb.NDF {
	file.RLock()
	defer file.RUnlock()

	return file.pb
}

// CompareHash evaluates if the passed NDF hash is the same as the stored one.
func (file *Ndf) CompareHash(h []byte) bool {
	file.RLock()
	defer file.RUnlock()

	// Return whether the hashes are different
	return bytes.Equal(file.hash, h)
}

// GenerateNDFHash generates a hash of the unmarshalled NDF bytes in the comms
// message. If the message or NDF bytes is nil, zeroes are returned.
func GenerateNDFHash(msg *pb.NDF) ([]byte, error) {
	if msg == nil || msg.GetNdf() == nil {
		return make([]byte, 32), nil
	}
	// Create new BLAKE2b hash
	hash, err := blake2b.New256(nil)
	if err != nil {
		return nil, err
	}

	hash.Write(msg.Ndf)

	return hash.Sum(nil), nil
}
