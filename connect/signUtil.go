package connect

import (
	"crypto/rand"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature/rsa"
)

func (c *ConnectionManager) SignMessage(anyMessage *any.Any, id string) (*pb.SignedMessage, error) {

	// Get hashed data
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	data := []byte(anyMessage.String())
	hashed := hash.Sum(data)[len(data):]

	key := c.GetPrivateKey()

	// Sign the thing
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		jww.ERROR.Printf("Failed to form message signature: %+v", err)
	}

	// Form signed message
	signedMessage := pb.SignedMessage{
		Message:   anyMessage,
		Signature: signature,
		ID:        id,
	}

	return &signedMessage, nil
}

func (c *ConnectionManager) VerifySignature(message *pb.SignedMessage, pb proto.Message, pubKey *rsa.PublicKey) error {
	err := ptypes.UnmarshalAny(message.Message, pb)
	if err != nil {
		jww.ERROR.Printf("Failed to unmarshal generic message, check your input message type: %+v", err)
		return err
	}

	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	s := pb.String()
	data := []byte(s)
	hashed := hash.Sum(data)[len(data):]

	err = rsa.Verify(pubKey, options.Hash, hashed, message.Signature, nil)

	return nil
}
