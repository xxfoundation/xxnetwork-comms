package network

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"sync"
)

// The Instance struct stores a combination of comms info and round info for servers
type Instance struct {
	comm *connect.ProtoComms

	partial      *SecuredNdf
	full         *SecuredNdf
	roundUpdates *ds.Updates
	roundData    *ds.Data

	roundlock sync.RWMutex
}

// Initializer for instance structs from base comms and NDF
func NewInstance(c *connect.ProtoComms, partial, full *ndf.NetworkDefinition) (*Instance, error) {
	partialNdf, err := NewSecuredNdf(partial)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not create secured partial ndf")
	}
	fullNdf, err := NewSecuredNdf(full)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not create secured full ndf")
	}
	return &Instance{
		c,
		partialNdf,
		fullNdf,
		&ds.Updates{},
		&ds.Data{},
		sync.RWMutex{},
	}, nil
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

// Return the partial ndf from this instance
func (i *Instance) GetPartialNdf() *SecuredNdf {
	return i.partial
}

// Return the full NDF from this instance
func (i *Instance) GetFullNdf() *SecuredNdf {
	return i.full
}

// Add a round to the round and update buffer
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

	i.roundlock.Lock()
	defer i.roundlock.Unlock()

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

// Get the round of a given ID
func (i *Instance) GetRound(id id.Round) (*pb.RoundInfo, error) {
	return i.roundData.GetRound(int(id))
}

// Get an update ID
func (i *Instance) GetRoundUpdate(updateID int) (*pb.RoundInfo, error) {
	return i.roundUpdates.GetUpdate(updateID)
}

// Get updates from a given round
func (i *Instance) GetRoundUpdates(id int) ([]*pb.RoundInfo, error) {
	return i.roundUpdates.GetUpdates(id)
}

// get the most recent update id
func (i *Instance) GetLastUpdateID() int {
	return i.roundUpdates.GetLastUpdateID()
}

// get the most recent round id
func (i *Instance) GetLastRoundID() id.Round {
	return i.roundData.GetLastRoundID()
}

// Update gateway hosts based on most complete ndf
func (i *Instance) UpdateGatewayConnections() error {
	if i.full != nil {
		return updateConns(i.full.f.Get(), i.comm, true, false)
	} else if i.partial != nil {
		return updateConns(i.partial.f.Get(), i.comm, true, false)
	} else {
		return errors.New("No ndf currently stored")
	}
}

// Update node hosts based on most complete ndf
func (i *Instance) UpdateNodeConnections() error {
	if i.full != nil {
		return updateConns(i.full.f.Get(), i.comm, false, true)
	} else if i.partial != nil {
		return updateConns(i.partial.f.Get(), i.comm, false, true)
	} else {
		return errors.New("No ndf currently stored")
	}
}

// Update host helper
func updateConns(def *ndf.NetworkDefinition, comms *connect.ProtoComms, gate, node bool) error {
	if gate {
		for i, h := range def.Gateways {
			gwid := id.NewNodeFromBytes(def.Nodes[i].ID).NewGateway().String()
			_, ok := comms.GetHost(gwid)
			if !ok {
				_, err := comms.AddHost(gwid, h.Address, []byte(h.TlsCertificate), false, true)
				if err != nil {
					return errors.WithMessagef(err, "Could not add gateway host %s", gwid)
				}
			}
		}
	}
	if node {
		for _, h := range def.Nodes {
			nid := id.NewNodeFromBytes(h.ID).String()
			_, ok := comms.GetHost(nid)
			if !ok {
				_, err := comms.AddHost(nid, h.Address, []byte(h.TlsCertificate), false, true)
				if err != nil {
					return errors.WithMessagef(err, "Could not add node host %s", nid)
				}
			}
		}
	}
	return nil
}
