package connect

import (
	"crypto/rand"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"golang.org/x/crypto/openpgp/errors"
)

// SignMessage takes a generic-type message and an ID, returns a SignedMessage
// The message is signed with the ConnectionManager's RSA PrivateKey
func (c *ConnectionManager) SignMessage(anyMessage *any.Any, id string) (*pb.SignedMessage, error) {
	// Get hashed data
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	data := []byte(anyMessage.String())
	hashed := hash.Sum(data)[len(data):]

	key := c.GetPrivateKey()
	if key == nil {
		return nil, errors.InvalidArgumentError("Connection manager: nil private key")
	}

	// Sign the thing
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		return nil, err
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
func (c *ConnectionManager) VerifySignature(message *pb.SignedMessage, pb proto.Message, pubKey *rsa.PublicKey) error {
	// Get hashed data of the message
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	s := pb.String()
	data := []byte(s)
	hashed := hash.Sum(data)[len(data):]

	// Verify signature of message
	err := rsa.Verify(pubKey, options.Hash, hashed, message.Signature, nil)
	if err != nil {
		return err
	}

	return nil
}
