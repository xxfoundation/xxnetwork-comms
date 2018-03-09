////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompShare
func TestSendPrecompShare(t *testing.T) {
	_, err := SendPrecompShare(SERVER_ADDRESS, &pb.PrecompShareMessage{})
	if err != nil {
		t.Errorf("PrecompShare: Error received: %s", err)
	}
}

// Smoke test SendPrecompShareInit
func TestSendPrecompShareInit(t *testing.T) {
	_, err := SendPrecompShareInit(SERVER_ADDRESS,
		&pb.PrecompShareInitMessage{})
	if err != nil {
		t.Errorf("PrecompShareInit: Error received: %s", err)
	}
}

// Smoke test SendPrecompShareCompare
func TestSendPrecompShareCompare(t *testing.T) {
	_, err := SendPrecompShareCompare(SERVER_ADDRESS,
		&pb.PrecompShareCompareMessage{})
	if err != nil {
		t.Errorf("PrecompShareCompare: Error received: %s", err)
	}
}

// Smoke test SendPrecompShareConfirm
func TestSendPrecompShareConfirm(t *testing.T) {
	_, err := SendPrecompShareConfirm(SERVER_ADDRESS,
		&pb.PrecompShareConfirmMessage{})
	if err != nil {
		t.Errorf("PrecompShareConfirm: Error received: %s", err)
	}
}

// Smoke test SendPrecompPermute
func TestSendPrecompPermute(t *testing.T) {
	_, err := SendPrecompPermute(SERVER_ADDRESS, &pb.PrecompPermuteMessage{})
	if err != nil {
		t.Errorf("PrecompPermute: Error received: %s", err)
	}
}

// Smoke test SendPrecompEncrypt
func TestSendPrecompEncrypt(t *testing.T) {
	_, err := SendPrecompEncrypt(SERVER_ADDRESS, &pb.PrecompEncryptMessage{})
	if err != nil {
		t.Errorf("PrecompEncrypt: Error received: %s", err)
	}
}

// Smoke test SendPrecompDecrypt
func TestSendPrecompDecrypt(t *testing.T) {
	_, err := SendPrecompDecrypt(SERVER_ADDRESS, &pb.PrecompDecryptMessage{})
	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
}

// Smoke test SendPrecompReveal
func TestSendPrecompReveal(t *testing.T) {
	_, err := SendPrecompReveal(SERVER_ADDRESS, &pb.PrecompRevealMessage{})
	if err != nil {
		t.Errorf("PrecompReveal: Error received: %s", err)
	}
}
