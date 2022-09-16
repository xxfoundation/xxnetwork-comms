////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handle basic logic for common operations of network instances

package network

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	ds "gitlab.com/elixxir/comms/network/dataStructures"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/ec"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/ndf"
	"testing"
)

// The Instance struct stores a combination of comms info and round info for servers
type Instance struct {
	comm            *connect.ProtoComms
	cmixGroup       *ds.Group // make a wrapper structure containing a group and a rwlock
	e2eGroup        *ds.Group
	partial         *SecuredNdf
	full            *SecuredNdf
	roundUpdates    *ds.Updates
	roundData       *ds.Data
	ers             ds.ExternalRoundStorage
	validationLevel ValidationType

	ipOverride *ds.IpOverrideList

	// Determines whether auth is enabled
	// on communication with gateways
	gatewayAuth bool

	// Network Health
	networkHealth chan Heartbeat

	// Determines whether verification for round info will be done
	// using the RSA key or the EC key.
	// Set to true, they shall use elliptic, set to false they shall use RSA
	useElliptic bool
	ecPublicKey *ec.PublicKey
	// Waiting Rounds
	waitingRounds *ds.WaitingRounds

	// Round Event Model
	events *ds.RoundEvents

	// Node Event Model Channels
	addNode       chan NodeGateway
	removeNode    chan *id.ID
	addGateway    chan NodeGateway
	removeGateway chan *id.ID
}

// Object used to signal information about the network health
type Heartbeat struct {
	HasWaitingRound bool
	IsRoundComplete bool
}

// Combines a node and gateway together together for return over channels
type NodeGateway struct {
	Node    ndf.Node
	Gateway ndf.Gateway
}

// Register NetworkHealth channel with Instance
func (i *Instance) SetNetworkHealthChan(c chan Heartbeat) {
	i.networkHealth = c
}

// Register AddNode channel with Instance
func (i *Instance) SetAddNodeChan(c chan NodeGateway) {
	i.addNode = c
}

// Register RemoveNode channel with Instance
func (i *Instance) SetRemoveNodeChan(c chan *id.ID) {
	i.removeNode = c
}

// Register AddGateway channel with Instance
func (i *Instance) SetAddGatewayChan(c chan NodeGateway) {
	i.addGateway = c
}

// Return AddGateway channel from Instance
func (i *Instance) GetAddGatewayChan() chan NodeGateway {
	return i.addGateway
}

// Register RemoveGateway channel with Instance
func (i *Instance) SetRemoveGatewayChan(c chan *id.ID) {
	i.removeGateway = c
}

// Return the Instance WaitingRounds object
func (i *Instance) GetWaitingRounds() *ds.WaitingRounds {
	return i.waitingRounds
}

// Return the Instance RoundEvents object
func (i *Instance) GetRoundEvents() *ds.RoundEvents {
	return i.events
}

// Return the partial ndf from this instance
func (i *Instance) GetPartialNdf() *SecuredNdf {
	return i.partial
}

// Return the full NDF from this instance
func (i *Instance) GetFullNdf() *SecuredNdf {
	return i.full
}

// Initializer for instance structs from base comms and NDF, you can put in nil for
// ERS if you don't want to use it
// useElliptic determines whether client will verify signatures using the RSA key
// or the elliptic curve key.
func NewInstance(c *connect.ProtoComms, partial, full *ndf.NetworkDefinition, ers ds.ExternalRoundStorage,
	validationLevel ValidationType, useElliptic bool) (*Instance, error) {
	var partialNdf *SecuredNdf
	var fullNdf *SecuredNdf
	var err error

	if partial == nil && full == nil {
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

		ipOverride:  ds.NewIpOverrideList(),
		useElliptic: useElliptic,
	}

	var ecPublicKey *ec.PublicKey
	if full != nil && full.Registration.EllipticPubKey != "" {
		ecPublicKey, err = ec.LoadPublicKey(i.GetEllipticPublicKey())
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Could not load elliptic key from ndf"))
		}
	} else if partial.Registration.EllipticPubKey != "" {
		ecPublicKey, err = ec.LoadPublicKey(i.GetEllipticPublicKey())
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Could not load elliptic key from ndf"))
		}
	}

	if ecPublicKey != nil {
		i.ecPublicKey = ecPublicKey
	} else {
		jww.DEBUG.Printf("Elliptic public key was not set, could not be found in NDF")
	}

	cmix := ""
	if full != nil && full.CMIX.Prime != "" {
		cmix, _ = full.CMIX.String()
	} else if partial.CMIX.Prime != "" {
		cmix, _ = partial.CMIX.String()
	}

	if cmix != "" {
		err := i.cmixGroup.Update(cmix)
		if err != nil {
			jww.WARN.Printf("Error updating cmix group: %+v", err)
		}
	}

	e2e := ""
	if full != nil && full.E2E.Prime != "" {
		e2e, _ = full.E2E.String()
	} else if partial.E2E.Prime != "" {
		e2e, _ = partial.E2E.String()
	}

	if cmix != "" {
		err := i.e2eGroup.Update(e2e)
		if err != nil {
			jww.WARN.Printf("Error updating e2e group: %+v", err)
		}
	}

	i.waitingRounds = ds.NewWaitingRounds()
	i.events = ds.NewRoundEvents()
	i.validationLevel = validationLevel

	// Set our ERS to the passed in ERS object (or nil)
	i.ers = ers

	return i, nil
}

// Utility function to create instance FOR TESTING PURPOSES ONLY
func NewInstanceTesting(c *connect.ProtoComms, partial, full *ndf.NetworkDefinition,
	e2eGroup, cmixGroup *cyclic.Group, i interface{}) (*Instance, error) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("NewInstanceTesting is restricted to testing only. Got %T", i)
	}
	instance, err := NewInstance(c, partial, full, nil, 0, false)
	if err != nil {
		return nil, errors.Errorf("Unable to create instance: %+v", err)
	}

	instance.cmixGroup.UpdateCyclicGroupTesting(cmixGroup, i)
	instance.e2eGroup.UpdateCyclicGroupTesting(e2eGroup, i)

	return instance, nil
}

//update the partial ndf
func (i *Instance) UpdatePartialNdf(m *pb.NDF) error {
	if i.partial == nil {
		return errors.New("Cannot update the partial ndf when it is nil")
	}

	perm, success := i.comm.GetHost(&id.Permissioning)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for NDF partial verification")
	}

	// Get a list of current nodes so we can check later for removed nodes
	oldNodeList := i.partial.Get().Nodes

	// Update the partial ndf
	err := i.partial.update(m, perm.GetPubKey())
	if err != nil {
		return err
	}

	// Get list of removed nodes and remove them from the host map
	rmNodes, err := getBannedNodes(oldNodeList, i.partial.Get().Nodes)
	if err != nil {
		return err
	}
	for _, nid := range rmNodes {
		i.comm.RemoveHost(nid)

		// Send events into Node Listener
		if i.removeNode != nil && i.removeGateway != nil {
			select {
			case i.removeNode <- nid:
			default:
				jww.WARN.Printf("Unable to send RemoveNode event for id %s", nid.String())
			}
			gwId := nid.DeepCopy()
			gwId.SetType(id.Gateway)
			select {
			case i.removeGateway <- nid:
			default:
				jww.WARN.Printf("Unable to send RemoveGateway event for id %s", nid.String())
			}
		}
	}

	// update the cmix group object
	cmixGrp, _ := i.partial.Get().CMIX.String()
	err = i.cmixGroup.Update(cmixGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update cmix group")
	}

	// update the e2e group object
	e2eGrp, _ := i.partial.Get().E2E.String()
	err = i.e2eGroup.Update(e2eGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update e2e group")
	}

	return nil
}

//overrides an IP address for an ID with one from
func (i *Instance) GetIpOverrideList() *ds.IpOverrideList {
	return i.ipOverride
}

//Gets the node and gateway with the given ID
func (i *Instance) GetNodeAndGateway(ngid *id.ID) (NodeGateway, error) {
	index := -1

	def := i.GetFullNdf()
	if def == nil {
		def = i.GetPartialNdf()
	}

	idBytes := ngid.Bytes()

	// depending on if the passed id is a node or gateway ID, look it up in the
	// correct list
	if ngid.GetType() == id.Node {
		for iter, n := range def.Get().Nodes {
			if bytes.Compare(n.ID, idBytes) == 0 {
				index = iter
				break
			}
		}
	} else if ngid.GetType() == id.Gateway {
		for iter, g := range def.Get().Gateways {
			if bytes.Compare(g.ID, idBytes) == 0 {
				index = iter
				break
			}
		}
	} else {
		return NodeGateway{}, errors.Errorf("The passed ID is not for "+
			"a node or gateway: %s", ngid)
	}

	//if no node or gateway is found, return an error
	if index == -1 {
		return NodeGateway{}, errors.Errorf("Failed to find Node or "+
			"Gateway with ID %s", ngid)
	}

	//return the found node and gateway
	return NodeGateway{
		Node:    def.Get().Nodes[index],
		Gateway: def.Get().Gateways[index],
	}, nil
}

//update the full ndf
func (i *Instance) UpdateFullNdf(m *pb.NDF) error {
	if i.full == nil {
		return errors.New("Cannot update the full ndf when it is nil")
	}

	perm, success := i.comm.GetHost(&id.Permissioning)

	if !success {
		return errors.New("Could not get permissioning Public Key" +
			"for full NDF verification")
	}

	// Get a list of current nodes so we can check later for removed nodes
	oldNodeList := i.full.Get().Nodes

	// Update the full ndf
	err := i.full.update(m, perm.GetPubKey())
	if err != nil {
		return err
	}
	rmNodes, err := getBannedNodes(oldNodeList, i.full.Get().Nodes)
	if err != nil {
		return err
	}
	for _, nid := range rmNodes {
		i.comm.RemoveHost(nid)

		// Send events into Node Listener
		if i.removeNode != nil {
			select {
			case i.removeNode <- nid:
			default:
				jww.WARN.Printf("Unable to send RemoveNode event for id %s", nid.String())
			}
		}
		// Send events into Gateway Listener
		if i.removeGateway != nil {
			gwId := nid.DeepCopy()
			gwId.SetType(id.Gateway)
			select {
			case i.removeGateway <- nid:
			default:
				jww.WARN.Printf("Unable to send RemoveGateway event for id %s", nid.String())
			}
		}
	}

	// update the cmix group object
	cmixGrp, _ := i.full.Get().CMIX.String()
	err = i.cmixGroup.Update(cmixGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update cmix group")
	}

	// update the e2e group object
	e2eGrp, _ := i.full.Get().E2E.String()
	err = i.e2eGroup.Update(e2eGrp)
	if err != nil {
		return errors.WithMessage(err, "Unable to update e2e group")
	}

	return nil

}

// Find nodes that have been removed, comparing two NDFs
func getBannedNodes(old []ndf.Node, new []ndf.Node) ([]*id.ID, error) {
	// List of nodes to get rid of
	var rmNodes []*id.ID
	// Get the list of old nodes and populate a map with them
	newNodeMap := make(map[string]ndf.Node)
	for _, n := range new {
		newNodeMap[string(n.ID)] = n
	}

	// Check the old nodes list against our new map
	for _, n := range old {
		// We try to find the "new" ID in the old map and see if it exists
		// If it doesn't exist, we remove it, since that means it's been removed from the network
		_, ok := newNodeMap[string(n.ID)]
		if !ok {
			nid, err := id.Unmarshal(n.ID)
			if err != nil {
				return nil, err
			}
			rmNodes = append(rmNodes, nid)
		}
	}

	return rmNodes, nil
}

// Pluralized version of RoundUpdate used by Client
func (i *Instance) RoundUpdates(rounds []*pb.RoundInfo) error {
	// Keep track of whether one of the rounds is completed
	isRoundComplete := false
	addedRounds := make([]*ds.Round, 0, len(rounds))
	removedRounds := make([]*ds.Round, 0, len(rounds))
	roundsToTrigger := make([]*ds.Round, 0, len(rounds))
	for _, round := range rounds {
		if states.Round(round.State) == states.COMPLETED {
			isRoundComplete = true
		}

		// Send the RoundUpdate
		rnd, err := i.RoundUpdate(round)
		if err != nil {
			return err
		}
		state := states.Round(round.State)
		if state == states.QUEUED {
			addedRounds = append(addedRounds, rnd)
		} else if state > states.QUEUED {
			addedRounds = append(removedRounds, rnd)
		}

		roundsToTrigger = append(roundsToTrigger, rnd)
	}

	go i.events.TriggerRoundEvents(roundsToTrigger...)

	i.waitingRounds.Insert(addedRounds, removedRounds)

	// Send a Heartbeat over the networkHealth channel
	if i.networkHealth != nil {
		select {
		case i.networkHealth <- Heartbeat{
			HasWaitingRound: i.GetWaitingRounds().Len() > 0,
			IsRoundComplete: isRoundComplete,
		}:
		default:
			jww.WARN.Printf("Unable to send NetworkHealth event")
		}
	}

	return nil
}

// Add a round to the round and update buffer
func (i *Instance) RoundUpdate(info *pb.RoundInfo) (*ds.Round, error) {
	perm, success := i.comm.GetHost(&id.Permissioning)

	if !success {
		return nil, errors.New("Could not get permissioning Public Key" +
			"for round info verification")
	}

	var rnd *ds.Round
	if i.useElliptic {
		// Use the elliptic key only
		rnd = ds.NewRound(info, nil, i.ecPublicKey)
	} else {
		// Use the rsa key only
		rnd = ds.NewRound(info, perm.GetPubKey(), nil)
	}

	if i.validationLevel == Strict {
		err := signature.VerifyRsa(info, perm.GetPubKey())
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Could not validate "+
				"the roundInfo signature: %+v", info))
		}
	}

	err := i.roundUpdates.AddRound(rnd)
	if err != nil {
		return nil, err
	}
	err = i.roundData.UpsertRound(rnd)
	if err != nil {
		return nil, err
	}
	if i.ers != nil {
		// If we are not lazy, we validate the info before storage
		if i.validationLevel != Lazy {
			_ = rnd.Get()
		}

		// Intentionally suppress error
		_ = i.ers.Store(info)
	}

	return rnd, nil
}

// GetE2EGroup gets the e2eGroup from the instance
func (i *Instance) GetE2EGroup() *cyclic.Group {
	return i.e2eGroup.Get()
}

// GetE2EGroup gets the cmixGroup from the instance
func (i *Instance) GetCmixGroup() *cyclic.Group {
	return i.cmixGroup.Get()
}

// Get the round of a given ID as a roundInfo (protobuff)
func (i *Instance) GetRound(id id.Round) (*pb.RoundInfo, error) {
	return i.roundData.GetRound(int(id))
}

// Get the round of a given ID as a ds.Round object
func (i *Instance) GetWrappedRound(id id.Round) (*ds.Round, error) {
	return i.roundData.GetWrappedRound(int(id))
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
	return i.roundData.GetLastRoundID() - 1
}

// Get the oldest round id
func (i *Instance) GetOldestRoundID() id.Round {
	return i.roundData.GetOldestRoundID()
}

// Update gateway hosts based on most complete ndf
func (i *Instance) UpdateGatewayConnections() error {
	if i.full != nil {
		return i.updateConns(i.full.f.Get(), true, false)
	} else if i.partial != nil {
		return i.updateConns(i.partial.f.Get(), true, false)
	} else {
		return errors.New("No ndf currently stored")
	}
}

// Update node hosts based on most complete ndf
func (i *Instance) UpdateNodeConnections() error {
	if i.full != nil {
		return i.updateConns(i.full.f.Get(), false, true)
	} else if i.partial != nil {
		return i.updateConns(i.partial.f.Get(), false, true)
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

// GetEllipticPublicKey gets the permissioning's elliptic public key
// from one of the NDFs
// It first checks the full ndf and returns if that has the key
// If not it checks the partial ndf and returns if it has it
// Otherwise it returns an empty string
func (i *Instance) GetEllipticPublicKey() string {
	// Check if the full ndf has the information
	if i.GetFullNdf() != nil {
		return i.GetFullNdf().Get().Registration.EllipticPubKey
	} else if i.GetPartialNdf() != nil {
		// Else check if the partial ndf has the information
		return i.GetPartialNdf().Get().Registration.EllipticPubKey
	}

	// If neither do, return an empty string
	return ""

}

// GetPermissioningId gets the permissioning ID from primitives
func (i *Instance) GetPermissioningId() *id.ID {
	return &id.Permissioning
}

// Update host helper
func (i *Instance) updateConns(def *ndf.NetworkDefinition, isGateway, isNode bool) error {
	if isGateway {
		for index, gateway := range def.Gateways {
			gwid, err := id.Unmarshal(def.Nodes[index].ID)
			if err != nil {
				return err
			}
			gwid.SetType(id.Gateway)
			//check if an ip override is registered
			addr := i.ipOverride.CheckOverride(gwid, gateway.Address)

			//check if the host exists
			host, ok := i.comm.GetHost(gwid)
			if !ok {

				// Check if gateway ID collides with an existing hard coded ID
				if id.CollidesWithHardCodedID(gwid) {
					return errors.Errorf("Gateway ID invalid, collides with a "+
						"hard coded ID. Invalid ID: %v", gwid.Marshal())
				}

				// If this entity is a gateway, other gateway hosts
				// should have auth enabled. Otherwise, disable auth
				gwParams := connect.GetDefaultHostParams()
				gwParams.MaxRetries = 3
				gwParams.EnableCoolOff = true
				gwParams.AuthEnabled = i.gatewayAuth
				_, err := i.comm.AddHost(gwid, addr, []byte(gateway.TlsCertificate), gwParams)
				if err != nil {
					return errors.WithMessagef(err, "Could not add gateway host %s", gwid)
				}

				// Send events into Node Listener
				if i.addGateway != nil {
					ng := NodeGateway{
						Node:    def.Nodes[index],
						Gateway: gateway,
					}

					select {
					case i.addGateway <- ng:
					default:
						jww.WARN.Printf("Unable to send AddGateway event for id %s", gwid.String())
					}
				}

			} else if host.GetAddress() != addr {
				host.UpdateAddress(addr)
			}
		}
	}
	if isNode {
		for index, node := range def.Nodes {
			nid, err := id.Unmarshal(node.ID)
			if err != nil {
				return err
			}
			//check if an ip override is registered
			addr := i.ipOverride.CheckOverride(nid, node.Address)

			//check if the host exists
			host, ok := i.comm.GetHost(nid)
			if !ok {

				// Check if isNode ID collides with an existing hard coded ID
				if id.CollidesWithHardCodedID(nid) {
					return errors.Errorf("Node ID invalid, collides with a "+
						"hard coded ID. Invalid ID: %v", nid.Marshal())
				}

				host, err := i.comm.AddHost(nid, addr, []byte(node.TlsCertificate), connect.GetDefaultHostParams())
				if err != nil {
					return errors.WithMessagef(err, "Could not add isNode host %s", nid)
				}

				// 10k batch size * 8192 packet size * 2
				host.SetWindowSize(connect.MaxWindowSize)

				// Send events into Node Listener
				if i.addNode != nil {
					ng := NodeGateway{
						Node:    node,
						Gateway: def.Gateways[index],
					}

					select {
					case i.addNode <- ng:
					default:
						jww.WARN.Printf("Unable to send AddNode event for id %s", nid.String())
					}
				}

			} else if host.GetAddress() != addr {
				host.UpdateAddress(addr)
			}
		}
	}
	return nil
}

// SetGatewayAuth will force authentication on all communications with gateways
// intended for use between Gateway <-> Gateway communications
func (i *Instance) SetGatewayAuthentication() {
	i.gatewayAuth = true
}
