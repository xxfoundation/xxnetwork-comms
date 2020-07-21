///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package mixmessages

import (
	"encoding/binary"
)

// GenerateSlotDigest serializes the gateway slot message for the
// client to hash
func GenerateSlotDigest(gatewaySlot *GatewaySlot) []byte {

	var gatewaySlotDigest []byte
	gatewaySlotDigest = append(gatewaySlotDigest, gatewaySlot.Message.SenderID...)
	gatewaySlotDigest = append(gatewaySlotDigest, gatewaySlot.Message.PayloadA...)
	gatewaySlotDigest = append(gatewaySlotDigest, gatewaySlot.Message.PayloadB...)

	for _, kmac := range gatewaySlot.Message.KMACs {
		gatewaySlotDigest = append(gatewaySlotDigest, kmac...)
	}

	roundIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(roundIdBytes, gatewaySlot.RoundID)

	gatewaySlotDigest = append(gatewaySlotDigest, roundIdBytes...)

	return gatewaySlotDigest

}
