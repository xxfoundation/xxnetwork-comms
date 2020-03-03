package network

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/primitives/id"
)

type Instance struct {
	comm *connect.ProtoComms

	partial      *SecuredNdf
	full         *SecuredNdf
	roundUpdates *ds.Updates
	roundData    *ds.Data
}

func NewInstance(c *connect.ProtoComms) *Instance {
	return &Instance{
		c,
		NewSecuredNdf(),
		NewSecuredNdf(),
		&ds.Updates{},
		&ds.Data{},
	}
}

//update the partial ndf
func (i *Instance) UpdatePartialNdf(m *pb.NDF) error {
	perm, success := i.comm.GetHost(id.PERMISSIONING)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for NDF partial verification")
	}

	return i.partial.update(m, perm.GetPubKey())
}

//update the full ndf
func (i *Instance) UpdateFullNdf(m *pb.NDF) error {
	perm, success := i.comm.GetHost(id.PERMISSIONING)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for full NDF verification")
	}

	return i.full.update(m, perm.GetPubKey())
}

func (i *Instance) GetPartialNdf() *SecuredNdf {
	return i.partial
}

func (i *Instance) GetFullNdf() *SecuredNdf {
	return i.full
}

func (i *Instance) RoundUpdate(info *pb.RoundInfo) error {
	perm, success := i.comm.GetHost(id.PERMISSIONING)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for round info verification")
	}

	err := signature.Verify(info, perm.GetPubKey())
	if err != nil {
		return errors.WithMessage(err, "Could not validate NDF")
	}

	err = i.roundUpdates.AddRound(info)
	if err != nil {
		return err
	}

	err = i.roundData.UpsertRound(info)
	if err != nil {
		return err
	}

	return nil
}

func (i *Instance) GetRound(id id.Round) (*pb.RoundInfo, error) {
	return i.roundData.GetRound(int(id))
}

func (i *Instance) GetRoundUpdate(updateID int) (*pb.RoundInfo, error) {
	return i.roundUpdates.GetUpdate(updateID)
}

func (i *Instance) GetRoundUpdates(id int) ([]*pb.RoundInfo, error) {
	return i.roundUpdates.GetUpdates(id)
}

func (i *Instance)GetLastUpdateID()int{
	return i.roundUpdates.GetLastUpdateID()
}

func (i *Instance)GetLastRoundID()id.Round{
	return i.roundData.GetLastRoundID()
}

