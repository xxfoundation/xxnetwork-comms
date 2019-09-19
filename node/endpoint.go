////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server gRPC endpoints

package node

// TODO: A lot of message types from gRPC are passed through, and a number of
//       errors that can occur are not accounted for.

import (
	"github.com/golang/protobuf/ptypes"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handle a Broadcasted Ask Online event
func (s *NodeComms) AskOnline(ctx context.Context, msg *pb.Ping) (
	*pb.Ack, error) {
	return &pb.Ack{}, nil
}

// DownloadTopology handles an incoming DownloadTopology event
func (s *NodeComms) DownloadTopology(ctx context.Context,
	msg *pb.SignedMessage) (*pb.Ack, error) {

	// fixme: this has got to be bad, we need to review this...
	go func() {
		// Get message sender ID
		sender := msg.ID
		pubKey := s.ConnectionManager.GetConnectionInfo(sender).RsaPublicKey

		// Unmarshal message to its original type
		original := pb.NodeTopology{}
		err := ptypes.UnmarshalAny(msg.Message, &original)
		if err != nil {
			jww.ERROR.Printf("Failed to unmarshal generic message, check your input message type: %+v", err)
			return
		}

		// Verify message contents
		var verified bool
		if pubKey != nil {
			err = s.ConnectionManager.VerifySignature(msg, &original, pubKey)
			if err != nil {
				jww.ERROR.Printf("Failed to verify message contents: %+v", err)
				return
			}
			verified = true
		} else {
			s := "WARNING: No public key found for connection, proceeding without signature verification"
			jww.WARN.Println(s)
			verified = false
		}

		senderAddress := s.ConnectionManager.GetConnectionInfo(msg.ID).Address
		ci := MessageInfo{
			Signature:      msg.Signature,
			ValidSignature: verified,
			Address:        senderAddress,
			SenderId:       msg.ID,
		}

		s.handler.DownloadTopology(&ci, &original)
	}()

	return &pb.Ack{}, nil
}

// Handle a NewRound event
func (s *NodeComms) CreateNewRound(ctx context.Context,
	msg *pb.RoundInfo) (*pb.Ack, error) {
	// Call the server handler to start a new round
	return &pb.Ack{}, s.handler.CreateNewRound(msg)
}

// PostNewBatch polls the first node and sends a batch when it is ready
func (s *NodeComms) PostNewBatch(ctx context.Context, msg *pb.Batch) (*pb.Ack, error) {
	// Call the server handler to post a new batch
	err := s.handler.PostNewBatch(msg)

	return &pb.Ack{}, err
}

// Handle a Phase event
func (s *NodeComms) PostPhase(ctx context.Context, msg *pb.Batch) (*pb.Ack,
	error) {
	// Call the server handler with the msg
	s.handler.PostPhase(msg)
	return &pb.Ack{}, nil
}

// Handle a phase event using a stream server
func (s *NodeComms) StreamPostPhase(server pb.Node_StreamPostPhaseServer) error {
	return s.handler.StreamPostPhase(server)
}

// Handle a PostRoundPublicKey message
func (s *NodeComms) PostRoundPublicKey(ctx context.Context,
	msg *pb.RoundPublicKey) (*pb.Ack, error) {
	// Call the server handler that receives the key share
	s.handler.PostRoundPublicKey(msg)
	return &pb.Ack{}, nil
}

// GetBufferInfo returns buffer size (number of completed precomputations)
func (s *NodeComms) GetRoundBufferInfo(ctx context.Context,
	msg *pb.RoundBufferInfo) (
	*pb.RoundBufferInfo, error) {
	bufSize, err := s.handler.GetRoundBufferInfo()
	if bufSize < 0 {
		bufSize = 0
	}
	size := uint32(bufSize)
	return &pb.RoundBufferInfo{RoundBufferSize: size}, err
}

// Handles Registration Nonce Communication
func (s *NodeComms) RequestNonce(ctx context.Context,
	msg *pb.NonceRequest) (*pb.Nonce, error) {

	// Obtain the nonce by passing to server
	nonce, pk, err := s.handler.RequestNonce(msg.GetSalt(),
		msg.GetClientRSAPubKey(), msg.GetClientDHPubKey(),
		msg.GetClientSignedByServer().Signature,
		msg.GetRequestSignature().Signature)

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the NonceMessage
	return &pb.Nonce{
		Nonce:    nonce,
		DHPubKey: pk,
		Error:    errMsg,
	}, err
}

// Handles Registration Nonce Confirmation
func (s *NodeComms) ConfirmRegistration(ctx context.Context,
	msg *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Obtain signed client public key by passing to server
	signature, err := s.handler.ConfirmRegistration(msg.GetUserID(), msg.NonceSignedByClient.Signature)

	// Obtain the error message, if any
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	// Return the RegistrationConfirmation
	return &pb.RegistrationConfirmation{
		ClientSignedByServer: &pb.RSASignature{
			Signature: signature,
		},
		Error: errMsg,
	}, err
}

// PostPrecompResult sends final Message and AD precomputations.
func (s *NodeComms) PostPrecompResult(ctx context.Context,
	msg *pb.Batch) (*pb.Ack, error) {
	// Call the server handler to start a new round
	err := s.handler.PostPrecompResult(msg.GetRound().GetID(),
		msg.GetSlots())
	return &pb.Ack{}, err
}

// FinishRealtime broadcasts to all nodes when the realtime is completed
func (s *NodeComms) FinishRealtime(ctx context.Context, msg *pb.RoundInfo) (*pb.Ack, error) {
	// Call the server handler to finish realtime
	err := s.handler.FinishRealtime(msg)

	return &pb.Ack{}, err
}

// GetCompletedBatch should return a completed batch that the calling gateway
// hasn't gotten before
func (s *NodeComms) GetCompletedBatch(ctx context.Context,
	msg *pb.Ping) (*pb.Batch, error) {
	return s.handler.GetCompletedBatch()
}

func (s *NodeComms) GetMeasure(ctx context.Context, msg *pb.RoundInfo) (*pb.RoundMetrics, error) {
	rm, err := s.handler.GetMeasure(msg)
	return rm, err
}

func (s *NodeComms) GetSignedCert(ctx context.Context, msg *pb.Ping) (*pb.SignedCerts, error) {
	rm, err := s.handler.GetSignedCert(msg)
	return rm, err
}

func (s *NodeComms) SendRoundTripPing(ping *pb.RoundTripPing) error {
	err := s.handler.SendRoundTripPing(ping)
	return err
}
