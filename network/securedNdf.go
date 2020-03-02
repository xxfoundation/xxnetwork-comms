package network

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/ndf"
)

// wraps the ndf data structure, expoting all the functionality expect the
// ability to change the ndf
type SecuredNdf struct {
	f *ds.Ndf
}

func NewSecuredNdf() *SecuredNdf {
	return &SecuredNdf{
		&ds.Ndf{},
	}
}

// unexported NDF update code
func (sndf *SecuredNdf) update(m *pb.NDF, key *rsa.PublicKey) error {
	err := signature.Verify(m, key)
	if err != nil {
		return errors.WithMessage(err, "Could not validate NDF")
	}

	return sndf.f.Update(m)
}

func (sndf *SecuredNdf) Get() *ndf.NetworkDefinition {
	return sndf.f.Get()
}

func (sndf *SecuredNdf) GetHash() []byte {
	return sndf.f.GetHash()
}

func (sndf *SecuredNdf) GetPb() *pb.NDF {
	return sndf.f.GetPb()
}

func (sndf *SecuredNdf) CompareHash(h []byte) (bool, error) {
	return sndf.f.CompareHash(h)
}
