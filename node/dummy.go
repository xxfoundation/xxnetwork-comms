////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Dummy implementation (so you can use for tests)
package node

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
)

// Blank struct implementing ServerHandler interface for testing purposes (Passing to StartServer)
type TestInterface struct{}

func (m TestInterface) NewRound(roundId string) {}

func (m TestInterface) RoundtripPing(message *pb.TimePing) {}

func (m TestInterface) ServerMetrics(message *pb.ServerMetricsMessage) {}

func (m TestInterface) SetPublicKey(roundId string, pkey []byte) {}

func (m TestInterface) PrecompDecrypt(message *pb.PrecompDecryptMessage) {}

func (m TestInterface) PrecompEncrypt(message *pb.PrecompEncryptMessage) {}

func (m TestInterface) PrecompReveal(message *pb.PrecompRevealMessage) {}

func (m TestInterface) PrecompPermute(message *pb.PrecompPermuteMessage) {}

func (m TestInterface) PrecompShare(message *pb.PrecompShareMessage) {}

func (m TestInterface) PrecompShareInit(message *pb.PrecompShareInitMessage) {}

func (m TestInterface) PrecompShareCompare(message *pb.
	PrecompShareCompareMessage) {
}

func (m TestInterface) PrecompShareConfirm(message *pb.
	PrecompShareConfirmMessage) {
}

func (m TestInterface) RealtimeDecrypt(message *pb.RealtimeDecryptMessage) {}

func (m TestInterface) RealtimeEncrypt(message *pb.RealtimeEncryptMessage) {}

func (m TestInterface) RealtimePermute(message *pb.RealtimePermuteMessage) {}

func (m TestInterface) ClientPoll(message *pb.ClientPollMessage) *pb.CmixMessage {
	return &pb.CmixMessage{}
}

func (m TestInterface) RequestContactList(message *pb.ContactPoll) *pb.
	ContactMessage {
	return &pb.ContactMessage{}
}

func (m TestInterface) UserUpsert(message *pb.UpsertUserMessage) {}

func (m TestInterface) PollRegistrationStatus(message *pb.
	RegistrationPoll) *pb.RegistrationConfirmation {
	return &pb.RegistrationConfirmation{}
}

func (m TestInterface) SetNick(message *pb.Contact) {}

func (m TestInterface) ReceiveMessageFromClient(message *pb.CmixMessage) {}

func (m TestInterface) StartRound(messages *pb.InputMessages) {}
