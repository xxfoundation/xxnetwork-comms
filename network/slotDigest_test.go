////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package network

import (
	"encoding/binary"
	"gitlab.com/elixxir/comms/mixmessages"
	"reflect"
	"testing"
)

// Tests that GenerateSlotDigest outputs a byte slice the length of the sum
// of its serialized components. Also checks if output matches precanned data
func TestGenerateSlotDigest(t *testing.T) {
	senderID := []byte("senderId")
	payloadA := []byte("payloadA")
	payloadB := []byte("payloadB")
	roundId := uint64(11420)
	kmacs := [][]byte{[]byte("kmac1"), []byte("kmac2")}

	// Craft message 1
	msg := &mixmessages.Slot{
		PayloadA: payloadA,
		PayloadB: payloadB,
		KMACs:    kmacs,
		SenderID: senderID,
	}

	gwSlot := &mixmessages.GatewaySlot{
		Message: msg,
		RoundID: roundId,
	}

	gwDigest := GenerateSlotDigest(gwSlot)

	roundIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(roundIdBytes, roundId)

	expectedLen := len(senderID) + len(payloadA) + len(payloadB) + len(roundIdBytes)
	for _, kmac := range kmacs {
		expectedLen += len(kmac)
	}

	if len(gwDigest) != expectedLen {
		t.Errorf("GenerateSlotDigest failed length test."+
			"\n\tExpected length: %v"+
			"\n\tReceived length: %v", expectedLen, len(gwDigest))
	}

	if !reflect.DeepEqual(gwDigest, precannedGatewayDigest) {
		t.Errorf("GenerateSlotDigest did not match expected."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", precannedGatewayDigest, gwDigest)
	}

}

// Test that GenerateSlotDigest generates the same output with the same input
func TestGenerateSlotDigest_Consistency(t *testing.T) {
	senderID := []byte("senderId")
	payloadA := []byte("payloadA")
	payloadB := []byte("payloadB")
	roundId := uint64(42)
	kmacs := [][]byte{[]byte("kmac1"), []byte("kmac2")}

	// Craft message 1
	msg := &mixmessages.Slot{
		PayloadA: payloadA,
		PayloadB: payloadB,
		KMACs:    kmacs,
		SenderID: senderID,
	}

	gwSlot := &mixmessages.GatewaySlot{
		Message: msg,
		RoundID: roundId,
	}

	gwDigest1 := GenerateSlotDigest(gwSlot)
	gwDigest2 := GenerateSlotDigest(gwSlot)

	if !reflect.DeepEqual(gwDigest1, gwDigest2) {
		t.Errorf("GenerateSlotDigest outputted different results with identical input."+
			"\n\tPrimary output: %v"+
			"\n\tSecondary output: %v", gwDigest1, gwDigest2)
	}

}

// Tests that GenerateSlotDigest produces different output with different input
func TestGenerateSlotDigest_Inconsistency(t *testing.T) {
	senderID := []byte("senderId")
	payloadA := []byte("payloadA")
	payloadB := []byte("payloadB")
	roundId := uint64(42)
	kmacs := [][]byte{[]byte("kmac1"), []byte("kmac2")}

	// Craft message 1
	msg1 := &mixmessages.Slot{
		PayloadA: payloadA,
		PayloadB: payloadB,
		KMACs:    kmacs,
		SenderID: senderID,
	}

	gwSlot1 := &mixmessages.GatewaySlot{
		Message: msg1,
		RoundID: roundId,
	}

	// Craft message 2 with swapped payloads
	msg2 := &mixmessages.Slot{
		PayloadA: payloadB,
		PayloadB: payloadA,
		KMACs:    kmacs,
		SenderID: senderID,
	}

	gwSlot2 := &mixmessages.GatewaySlot{
		Message: msg2,
		RoundID: roundId,
	}

	// Generate slot digest
	gwDigest1 := GenerateSlotDigest(gwSlot1)
	gwDigest2 := GenerateSlotDigest(gwSlot2)

	if reflect.DeepEqual(gwDigest1, gwDigest2) {
		t.Errorf("GenerateSlotDigest outputted identical results with different input."+
			"\n\tPrimary output: %v"+
			"\n\tSecondary output: %v", gwDigest1, gwDigest2)
	}
}

var precannedGatewayDigest = []byte{115, 101, 110, 100, 101, 114, 73, 100, 112, 97, 121, 108, 111, 97, 100, 65, 112, 97, 121, 108, 111, 97, 100, 66, 107, 109, 97, 99, 49, 107, 109, 97, 99, 50, 0, 0, 0, 0, 0, 0, 44, 156}
