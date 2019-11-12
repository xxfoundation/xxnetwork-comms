package connect

import (
	"crypto/rand"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature/rsa"
)

// SignMessage takes a generic-type message and an ID, returns a SignedMessage
// The message is signed with the Manager's RSA PrivateKey
func (m *Manager) SignMessage(anyMessage *any.Any, id string) (*pb.SignedMessage, error) {
	// Get hashed data
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	data := []byte(anyMessage.String())
	hashed := hash.Sum(data)[len(data):]

	key := m.GetPrivateKey()
	if key == nil {
		jwalterweatherman.WARN.Printf("Private key was nil, sending message unsigned")
		return &pb.SignedMessage{
			Message:   anyMessage,
			Signature: nil,
			ID:        id,
		}, nil
	}

	// Sign the thing
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Form signed message
	signedMessage := pb.SignedMessage{
		Message:   anyMessage,
		Signature: signature,
		ID:        id,
	}

	return &signedMessage, nil
}

// VerifySignature accepts a signed message, the UnMarshalled message, and RSA PublicKey
// It verifies the signature, returning an error if invalid
func (m *Manager) VerifySignature(message *pb.SignedMessage,
	pb proto.Message, host *Host) error {

	// Get hashed data of the message
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	s := pb.String()
	data := []byte(s)
	hashed := hash.Sum(data)[len(data):]

	// Verify signature of message
	err := rsa.Verify(host.rsaPublicKey, options.Hash, hashed, message.Signature, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
