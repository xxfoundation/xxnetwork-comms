////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handle basic logic for common operations of network instances

package network

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"testing"
)

// The Instance struct stores a combination of comms info and round info for servers
type Instance struct {
	comm         *connect.ProtoComms
	cmixGroup    *ds.Group // make a wrapper structure containing a group and a rwlock
	e2eGroup     *ds.Group
	partial      *SecuredNdf
	full         *SecuredNdf
	roundUpdates *ds.Updates
	roundData    *ds.Data
}

// Initializer for instance structs from base comms and NDF
func NewInstance(c *connect.ProtoComms, partial, full *ndf.NetworkDefinition) (*Instance, error) {
	var partialNdf *SecuredNdf
	var fullNdf *SecuredNdf
	var err error

	if partial == nil && full == nil{
		return nil, errors.New("Cannot create a network instance without an NDF")
	}

	if partial != nil {
		partialNdf, err = NewSecuredNdf(partial)
		if err != nil {
			return nil, errors.WithMessage(err, "Could not create secured partial ndf")
		}
	}

	if full != nil {
		fullNdf, err = NewSecuredNdf(full)
		if err != nil {
			return nil, errors.WithMessage(err, "Could not create secured full ndf")
		}
	}

	i := &Instance{
		comm:         c,
		partial:      partialNdf,
		full:         fullNdf,
		roundUpdates: ds.NewUpdates(),
		roundData:    ds.NewData(),
		cmixGroup:    ds.NewGroup(),
		e2eGroup:     ds.NewGroup(),
	}

	cmix, _ := partial.CMIX.String()
	if full.CMIX.Prime != "" {
		cmix, _ = full.CMIX.String()
	}
	if cmix == "" {
		return nil, errors.New("No cmix group was found in either NDF")
	} else {
		_ = i.cmixGroup.Update(cmix)
	}

	e2e, _ := partial.E2E.String()
	if full.E2E.Prime != "" {
		e2e, _ = full.E2E.String()
	}
	if cmix == "" {
		return nil, errors.New("No E2E group was found in either NDF")
	} else {
		_ = i.e2eGroup.Update(e2e)
	}

	return i, nil
}

// Utility function to create instance FOR TESTING PURPOSES ONLY
func NewInstanceTesting(c *connect.ProtoComms, partial, full *ndf.NetworkDefinition,
	e2eGroup, cmixGroup *cyclic.Group, t *testing.T) (*Instance, error) {
	if t == nil {
		panic("This is a utility function for testing purposes only!")
	}
	instance, err := NewInstance(c, partial, full)
	if err != nil {
		return nil, errors.Errorf("Unable to create instance: %+v", err)
	}

	instance.cmixGroup.UpdateCyclicGroupTesting(cmixGroup, t)
	instance.e2eGroup.UpdateCyclicGroupTesting(e2eGroup, t)

	return instance, nil
}

//update the partial ndf
func (i *Instance) UpdatePartialNdf(m *pb.NDF) error {
	if i.partial ==nil{
		return errors.New("Cannot update the partial ndf when it is nil")
	}

	perm, success := i.comm.GetHost(id.PERMISSIONING)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for NDF partial verification")
	}

	// Update the partial ndf
	err := i.partial.update(m, perm.GetPubKey())
	if err != nil {
		return err
	}

	// update the cmix group object
	cmixGrp, _ := i.partial.Get().CMIX.String()
	err = i.cmixGroup.Update(cmixGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update cmix group")
	}

	// update the cmix group object
	e2eGrp, _ := i.partial.Get().E2E.String()
	err = i.cmixGroup.Update(e2eGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update e2e group")
	}

	return nil

}

//update the full ndf
func (i *Instance) UpdateFullNdf(m *pb.NDF) error {
	if i.full ==nil{
		return errors.New("Cannot update the full ndf when it is nil")
	}

	perm, success := i.comm.GetHost(id.PERMISSIONING)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for full NDF verification")
	}

	// Update the full ndf
	err := i.full.update(m, perm.GetPubKey())
	if err != nil {
		return err
	}

	// update the cmix group object
	cmixGrp, _ := i.full.Get().CMIX.String()
	err = i.cmixGroup.Update(cmixGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update cmix group")
	}

	// update the cmix group object
	e2eGrp, _ := i.full.Get().E2E.String()
	err = i.cmixGroup.Update(e2eGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update e2e group")
	}

	return nil

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

// GetE2EGroup gets the e2eGroup from the instance
func (i *Instance) GetE2EGroup() *cyclic.Group {
	return i.e2eGroup.Get()
}

// GetE2EGroup gets the cmixGroup from the instance
func (i *Instance) GetCmixGroup() *cyclic.Group {

	return i.cmixGroup.Get()
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
func (i *Instance) GetRoundUpdates(id int) []*pb.RoundInfo {
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

// GetPermissioningAddress gets the permissioning address from one of the NDF
// It first checks the full ndf and returns if that has the address
// If not it checks the partial ndf and returns if it has it
// Otherwise it returns an empty string
func (i *Instance) GetPermissioningAddress() string {
	// Check if the full ndf has the information
	if i.GetFullNdf() != nil {
		return i.GetFullNdf().Get().Registration.Address
	} else if i.GetPartialNdf() != nil {
		// Else check if the partial ndf has the information
		return i.GetPartialNdf().Get().Registration.Address
	}

	// If neither do, return an empty string
	return ""
}

// GetPermissioningCert gets the permissioning certificate from one of the NDFs
// It first checks the full ndf and returns if that has the cert
// If not it checks the partial ndf and returns if it has it
// Otherwise it returns an empty string
func (i *Instance) GetPermissioningCert() string {
	// Check if the full ndf has the information
	if i.GetFullNdf() != nil {
		return i.GetFullNdf().Get().Registration.TlsCertificate
	} else if i.GetPartialNdf() != nil {
		// Else check if the partial ndf has the information
		return i.GetPartialNdf().Get().Registration.TlsCertificate
	}

	// If neither do, return an empty string
	return ""

}

// GetPermissioningId gets the permissioning ID from primitives
func (i *Instance) GetPermissioningId() string {
	return id.PERMISSIONING

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
