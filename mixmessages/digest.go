package mixmessages

import (
	"github.com/golang/protobuf/proto"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/hash"
)

// Function to digest Identity
func (i *Identity) Digest() []byte {
	// return hash(username|dhPubKey|salt)}
	// Generate the hash function
	h, err := hash.NewCMixHash()
	if err != nil {
		jww.FATAL.Panicf("Could not get hash: %+v", err)
	}

	// Marshal the message to put into the hash
	mb, err := proto.Marshal(i)
	if err != nil {
		jww.FATAL.Panicf("Could not marshal: %+v", err)
	}

	// Hash the Identity data to generate the vector
	h.Write(mb)
	return h.Sum(nil)
}

// Function to digest Fact
func (i *Fact) Digest() []byte {
	//return hash(Fact |FactType )}
	// Generate the hash function
	h, err := hash.NewCMixHash()
	if err != nil {
		jww.FATAL.Panicf("Could not get hash: %+v", err)
	}

	// Marshal the message to put into the hash
	mb, err := proto.Marshal(i)
	if err != nil {
		jww.FATAL.Panicf("Could not marshal: %+v", err)
	}

	// Hash the Fact data to generate the vector
	h.Write(mb)
	return h.Sum(nil)
}
