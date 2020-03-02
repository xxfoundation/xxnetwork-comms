package dataStructures

import (
	"bytes"
	"crypto/sha256"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	id2 "gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"sync"
)

type Ndf struct {
	f    *ndf.NetworkDefinition
	pb   *pb.NDF
	hash []byte
	lock sync.RWMutex
}

//Updates to a new NDF if the passed NDF is valid
func (file *Ndf) Update(m *pb.NDF, comm *connect.ProtoComms) error {

	//build the ndf object
	decoded, _, err := ndf.DecodeNDF(string(m.Ndf))

	if err != nil {
		return errors.WithMessage(err, "Could not decode the NDF")
	}

	file.lock.Lock()
	defer file.lock.Unlock()

	file.pb = m
	file.f = decoded

	for _, n := range file.f.Nodes {
		id := id2.NewNodeFromBytes(n.ID).String()
		_, ok := comm.GetHost(id)
		if !ok {
			_, err := comm.AddHost(id, n.Address, []byte(n.TlsCertificate), false, true)
			if err != nil {
				return errors.Wrapf(err, "Error adding host for node %s", id)
			}
		}
	}

	//set the ndf hash
	marshaled, err := file.f.Marshal()
	if err != nil {
		return errors.WithMessage(err,
			"Could not marshal NDF for hashing")
	}

	// Serialize then hash the constructed ndf
	hash := sha256.New()
	hash.Write(marshaled)
	file.hash = hash.Sum(nil)

	return nil
}

//returns the ndf object
//fix-me: return a copy instead to ensure edits to not impact the
//original version
func (file *Ndf) Get() *ndf.NetworkDefinition {
	file.lock.RLock()
	defer file.lock.RUnlock()

	return file.f
}

//returns the ndf hash
func (file *Ndf) GetHash() []byte {
	file.lock.RLock()
	defer file.lock.RUnlock()

	rtn := make([]byte, len(file.hash))
	copy(rtn, file.hash)
	return rtn
}

//returns the ndf hash
//fix-me: return a copy instead to ensure edits to not impact the
//original version
func (file *Ndf) GetPb() *pb.NDF {
	file.lock.RLock()
	defer file.lock.RUnlock()

	return file.pb
}

//evaluates if the passed ndf hash is the same as the stored one
//returns an error if no ndf is available, returns false if they are different
//and true if they are the same
func (file *Ndf) CompareHash(h []byte) (bool, error) {
	file.lock.RLock()
	defer file.lock.RUnlock()

	//return the NO_NDF error if no NDF is available
	if len(file.hash) == 0 {
		errMsg := errors.Errorf(ndf.NO_NDF)
		return false, errMsg
	}

	//return true if the hashes are the same
	if bytes.Compare(file.hash, h) == 0 {
		return true, nil
	}

	//return false if the hashes are different
	return false, nil
}
