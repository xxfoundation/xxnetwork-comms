///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"github.com/katzenpost/core/crypto/eddsa"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"sync/atomic"
)

// Structure wraps a round info object with the
// key to verify the protobuff's signature
// and a state track for verifying
type Round struct {
	info            *pb.RoundInfo
	needsValidation *uint32
	rsaPubKey       *rsa.PublicKey
	ecPubKey        *eddsa.PublicKey
}

// Constructor of a Round object.
func NewRound(ri *pb.RoundInfo, pubkey *rsa.PublicKey, ecPubKey *eddsa.PublicKey) *Round {
	validationDefault := uint32(0)
	return &Round{
		info:            ri,
		needsValidation: &validationDefault,
		rsaPubKey:       pubkey,
		ecPubKey:        ecPubKey,
	}
}

// Constructor of an already verified round object
// Intended for use by round creator.
func NewVerifiedRound(ri *pb.RoundInfo, pubkey *rsa.PublicKey) *Round {
	// Set validation to done on creation
	validationDefault := uint32(1)
	return &Round{
		info:            ri,
		needsValidation: &validationDefault,
		rsaPubKey:       pubkey,
	}
}

// Get returns the round info object. If we have not
// validated the signature before, we then verify.
// Later calls will not need validation
func (r *Round) Get() *pb.RoundInfo {
	if atomic.LoadUint32(r.needsValidation) == 0 {
		if r.rsaPubKey != nil {
			// Check the sig, panic if failure
			err := signature.Verify(r.info, r.rsaPubKey)
			if err != nil {
				jww.FATAL.Panicf("Could not validate "+
					"the roundInfo signature: %+v: %v", r.info, err)
			}
		}

		if r.ecPubKey != nil {
			// Check the sig, panic if failure
			err := signature.VerifyEddsa(r.info, r.ecPubKey)
			if err != nil {
				jww.FATAL.Panicf("Could not validate "+
					"the roundInfo signature: %+v: %v", r.info, err)
			}
		}

		atomic.StoreUint32(r.needsValidation, 1)
	}
	return r.info
}
