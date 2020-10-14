package mixmessages

import (
	"encoding/binary"
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

	// Hash the auth key to generate the vector
	h.Write([]byte(i.Username))
	h.Write(i.DhPubKey)
	h.Write(i.Salt)
	authKeyHash := h.Sum(nil)
	return authKeyHash
}

// Function to digest Fact
func (i *Fact) Digest() []byte {
	//return hash(Fact |FactType )}
	// Generate the hash function
	h, err := hash.NewCMixHash()
	if err != nil {
		jww.FATAL.Panicf("Could not get hash: %+v", err)
	}

	// Hash the auth key to generate the vector
	h.Write([]byte(i.Fact))

	// Convert FactType uint32 to []byte and write it to hash
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, i.FactType)
	h.Write(bs)

	authKeyHash := h.Sum(nil)
	return authKeyHash
}

// Function to digest FactRemoval
func (fr *FactRemoval) Digest() []byte {
	//return hash(Fact |FactType )
	// Generate the hash function
	h, err := hash.NewCMixHash()
	if err != nil {
		jww.FATAL.Panicf("Could not get hash: %+v", err)
	}

	// Hash the auth key to generate the vector
	h.Write([]byte(fr.Fact))

	// Convert FactType uint32 to []byte and write it to hash
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, fr.FactType)
	h.Write(bs)

	authKeyHash := h.Sum(nil)
	return authKeyHash
}
