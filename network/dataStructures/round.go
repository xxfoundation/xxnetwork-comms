///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"sync/atomic"
)

type Round struct{
	info            *pb.RoundInfo
	needsValidation *uint32
	pubkey          *rsa.PublicKey
}

func NewRound(ri *pb.RoundInfo, pubkey *rsa.PublicKey) *Round {
	validationDefault := uint32(0)
	return &Round{
		info:            ri,
		needsValidation: &validationDefault,
		pubkey:          pubkey,
	}
}

func (r *Round) Get() *pb.RoundInfo{
	if atomic.LoadUint32(r.needsValidation) == 0 {
		// Check the sig, panic if failure
		err := signature.Verify(r.info, r.pubkey)
		if err != nil {
			jww.FATAL.Panicf("Could not validate "+
				"the roundInfo signature: %+v: %v", r.info, err)
		}

		atomic.StoreUint32(r.needsValidation,1)
	}
	return r.info
}