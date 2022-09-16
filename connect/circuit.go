////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
)

type Circuit struct {
	nodes       []*id.ID
	nodeIndexes map[id.ID]int
	hosts       []*Host
}

// New makes a list of node addresses for use.  It finds
// the passed "myId" and denotes it internally for use with
// utility functions.  The nodeID are copied instead of linked
// to ensure any modification of them does not change the
// Circuit structure.  Will panic if the length of the passed
// list is zero.
func NewCircuit(list []*id.ID) *Circuit {
	c := Circuit{
		nodes:       make([]*id.ID, 0),
		nodeIndexes: make(map[id.ID]int),
		hosts:       make([]*Host, 0),
	}

	if len(list) == 0 {
		jww.FATAL.Panicf("Cannot build a Circuit of len 0")
	}

	for index, nid := range list {
		if _, ok := c.nodeIndexes[*nid]; ok {
			jww.FATAL.Panicf("NodeIDs must be unique for the circuit.Circuit, "+
				"%s passed multiple times", nid)
		}

		c.nodeIndexes[*nid] = index
		c.nodes = append(c.nodes, nid.DeepCopy())
	}

	return &c
}

// GetNodeLocation returns the location of the passed node in the list.
// Returns -1 if the node is not in the list
func (c *Circuit) GetNodeLocation(node *id.ID) int {
	nodeLoc, ok := c.nodeIndexes[*node]

	if !ok {
		return -1
	}

	return nodeLoc
}

// GetNodeAtIndex returns the node at the given index.  Panics
// if the index does not exist within the circuit
func (c *Circuit) GetNodeAtIndex(index int) *id.ID {
	if index < 0 || index >= len(c.nodes) {
		jww.FATAL.Panicf("Cannot get an index %v which is outside"+
			" the Circut (len=%v)", index, len(c.nodes))
	}
	return c.nodes[index].DeepCopy()
}

// Get the last node in the circuit, will panic if the circuit has nil as a node
func (c *Circuit) GetLastNode() *id.ID {
	return c.GetNodeAtIndex(c.Len() - 1)
}

// Len returns the number of nodes in the circuit
func (c *Circuit) Len() int {
	return len(c.nodes)
}

// GetNextNode gets the node following the passed node in
// the list. It wraps around to the beginning of the list
// if the passed node is the last node.
func (c *Circuit) GetNextNode(from *id.ID) *id.ID {
	loc := c.GetNodeLocation(from)

	if loc == -1 {
		jww.FATAL.Panicf("Cannot get the next node in the circuit.Circut"+
			"for node %s which is not present", from)
	}

	return c.nodes[(loc+1)%len(c.nodes)].DeepCopy()
}

// GetNextNode gets the node preceding the passed node in
// the list. It wraps around to the end of the list
// if the passed node is the first node.
func (c *Circuit) GetPrevNode(from *id.ID) *id.ID {
	loc := c.GetNodeLocation(from)

	if loc == -1 {
		jww.FATAL.Panicf("Cannot get the previous node in the circuit.Circut"+
			"for node %s which is not present", from)
	}

	var prevLoc int
	if loc == 0 {
		prevLoc = len(c.nodes) - 1
	} else {
		prevLoc = loc - 1
	}

	return c.nodes[prevLoc].DeepCopy()
}

// IsFirstNode returns true if the passed node is the
// first node, false otherwise
func (c *Circuit) IsFirstNode(node *id.ID) bool {
	return c.GetNodeLocation(node) == 0
}

// IsLastNode returns true if the passed node is the
// last node, false otherwise
func (c *Circuit) IsLastNode(node *id.ID) bool {
	return c.GetNodeLocation(node) == c.Len()-1
}

// GetOrdering returns a slice of Circuits with each one having a different
// shifted ordering.
func (c *Circuit) GetOrdering() []*Circuit {
	list := c.nodes
	circuits := make([]*Circuit, len(list))

	for i := range list {
		circuits[i] = NewCircuit(shiftLeft(list, i))
	}

	return circuits
}

//GetHostAtIndex: Gets host at requested index. Panics if index is outside
// of the range of the list
func (c *Circuit) GetHostAtIndex(index int) *Host {
	if index < 0 || index >= len(c.hosts) {
		jww.FATAL.Panicf("Cannot get an index %v which is outside"+
			" the Circut (len=%v)", index, len(c.hosts))
	}
	return c.hosts[index]
}

//SetHosts takes a list of hosts and copies them into the list of hosts in
// the circuit object
func (c *Circuit) AddHost(newHost *Host) {
	c.hosts = append(c.hosts, newHost)
}

// shiftLeft rotates the node IDs in a slice to the left the specified number of
// times.
func shiftLeft(list []*id.ID, rotation int) []*id.ID {
	var newList []*id.ID
	size := len(list)

	for i := 0; i < rotation; i++ {
		newList = list[1:size]
		newList = append(newList, list[0])
		list = newList
	}

	return list
}
