////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompShare
func TestSendPrecompShare(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompShare(ServerAddress, &pb.PrecompShareMessage{})
	if err != nil {
		t.Errorf("PrecompShare: Error received: %s", err)
	}
}

// Smoke test SendPrecompShareInit
func TestSendPrecompShareInit(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompShareInit(ServerAddress,
		&pb.PrecompShareInitMessage{})
	if err != nil {
		t.Errorf("PrecompShareInit: Error received: %s", err)
	}
}

// Smoke test SendPrecompShareCompare
func TestSendPrecompShareCompare(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompShareCompare(ServerAddress,
		&pb.PrecompShareCompareMessage{})
	if err != nil {
		t.Errorf("PrecompShareCompare: Error received: %s", err)
	}
}

// Smoke test SendPrecompShareConfirm
func TestSendPrecompShareConfirm(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompShareConfirm(ServerAddress,
		&pb.PrecompShareConfirmMessage{})
	if err != nil {
		t.Errorf("PrecompShareConfirm: Error received: %s", err)
	}
}

// Smoke test SendPrecompPermute
func TestSendPrecompPermute(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompPermute(ServerAddress, &pb.PrecompPermuteMessage{})
	if err != nil {
		t.Errorf("PrecompPermute: Error received: %s", err)
	}
}

// Smoke test SendPrecompEncrypt
func TestSendPrecompEncrypt(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompEncrypt(ServerAddress, &pb.PrecompEncryptMessage{})
	if err != nil {
		t.Errorf("PrecompEncrypt: Error received: %s", err)
	}
}

// Smoke test SendPrecompDecrypt
func TestSendPrecompDecrypt(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompDecrypt(ServerAddress, &pb.PrecompDecryptMessage{})
	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
}

// Smoke test SendPrecompReveal
func TestSendPrecompReveal(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPrecompReveal(ServerAddress, &pb.PrecompRevealMessage{})
	if err != nil {
		t.Errorf("PrecompReveal: Error received: %s", err)
	}
}
