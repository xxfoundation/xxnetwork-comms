////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"reflect"
	"testing"
)

//Tests the happy path of NewCircuit
func TestNew(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	//check that the internal id list matches the passed one
	same := true
	for index, nid := range nodeIdList {
		same = same && nid.Cmp(circuit.nodes[index])
	}

	if !same {
		t.Errorf("Circuit: internal list does not the same as the passed list")
	}

	//check that the indexes in the map are correct
	for index, nid := range nodeIdList {
		if circuit.nodeIndexes[*nid] != index {
			t.Errorf("Circuit: index linkage of %s incorrect; "+
				"Expected %v, Received: %v", nid, index,
				circuit.nodeIndexes[*nid])
		}
	}

	// check that the indexes in the map represent locations in the list
	// which contain the same node id which was used in the initial lookup
	for index, nid := range nodeIdList {
		if !reflect.DeepEqual(nid, circuit.nodes[circuit.nodeIndexes[*nid]]) {
			t.Errorf("Circuit: a index %v linkage of %s mismatch; "+
				"Expected %s, Received: %s", index, nid, nid,
				circuit.nodes[circuit.nodeIndexes[*nid]])
		}
	}

	//check that the internal list's internal data is not linked
	for index := range nodeIdList {
		nodeIdList[index][2] = 5
	}
	if reflect.DeepEqual(nodeIdList, circuit.nodes) {
		t.Errorf("Circuit: internal list linked to passed list")
	}
}

//Tests that New circuit properly errors when a list with duplicate nodes is passed
func TestNew_Duplicate(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	nodeIdList = append(nodeIdList, nodeIdList[1].DeepCopy())

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	NewCircuit(nodeIdList)

	t.Errorf("Circuit: no error when list contains duplicate node")

}

//Tests that New circuit properly errors when a list with duplicate nodes is passed
func TestNew_LenZero(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	NewCircuit(make([]*id.ID, 0))

	t.Errorf("Circuit: no error when creating with list of length zero")

}

// Tests the GetNodeLocation returns the correct location for all
// present nodeIDs
func TestCircuit_GetNodeLocation(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	for index, nid := range nodeIdList {
		if index != circuit.GetNodeLocation(nid) {
			t.Errorf("Circuit.GetNodeLocation: node location for node %s incorrect;"+
				"Expected: %v, Received: %v", nid, index, circuit.GetNodeLocation(nid))
		}
	}
}

// Tests the GetNodeLocation returns -1 when the node is not present
// present nodeIDs
func TestCircuit_GetNodeLocation_OutOfList(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	invalidNodeID := makeNodeId(77, t)

	invalidLoc := circuit.GetNodeLocation(invalidNodeID)
	if invalidLoc != -1 {
		t.Errorf("Circuit.GetNodeLocation: location returned when passed id (%s) is invalid:"+
			"Expected: -1, Received: %v", invalidNodeID, invalidLoc)
	}
}

// Tests the happy path of GetNodeAtIndex
func TestCircuit_GetNodeAtIndex(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	for index, nid := range nodeIdList {
		if !nid.Cmp(circuit.GetNodeAtIndex(index)) {
			t.Errorf("Circuit.GetNodeAtIndex: node at index %v incorrect;"+
				"Expected: %v, Received: %v", index, nid, circuit.GetNodeAtIndex(index))
		}
	}
}

// Tests that GetNodeAtIndex panics when the passed
// index is lower than 0
func TestCircuit_GetNodeAtIndex_OutOfList_Lower(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	circuit.GetNodeAtIndex(-1)

	t.Errorf("Circuit.GetNodeAtIndex: should have paniced with index of -1")

}

// Tests that GetNodeAtIndex panics when the passed
// index is len or greater
func TestCircuit_GetNodeAtIndex_OutOfList_Greater(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	circuit.GetNodeAtIndex(len(circuit.nodes))

	t.Errorf("Circuit.GetNodeAtIndex: should have paniced with index of len()")

}

//Tests that GetHostAtIndex panics whn the passed index
// is len or greater
func TestCircuit_GetHostAtIndex_OutOfList_Greater(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	circuit.GetHostAtIndex(len(nodeIdList))

	t.Errorf("Circuit.GetHostAtIndex: should have paniced with index of len()")

}

// Tests that GetNodeAtIndex panics when the passed
// index is lower than 0
func TestCircuit_GetHostAtIndex_OutOfList_Lesser(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)

	circuit := NewCircuit(nodeIdList)

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	circuit.GetNodeAtIndex(-1)

	t.Errorf("Circuit.GetHostAtIndex: should have paniced with index of -1")

}

// Tests the happy path of GetHostAtIndex
func TestCircuit_GetHostAtIndex(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(5, t)
	cert, _ := utils.ReadFile(testkeys.GetNodeCertPath())
	circuit := NewCircuit(nodeIdList)
	testID := id.NewIdFromString("test", id.Generic, t)
	testHost, _ := NewHost(testID, "test", cert, GetDefaultHostParams())
	circuit.AddHost(testHost)

	if !reflect.DeepEqual(circuit.hosts[0], testHost) {
		t.Errorf("Circuit.GetHostAtIndex: host incorrect;\n"+
			"Expected: %v\nReceived: %v", testHost, circuit.hosts[0])
	}

}

//Tests to see if node retrieved is in fact the last node
func TestCircuit_GetLastNode(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(23, t)
	circuit := NewCircuit(nodeIdList)
	if !circuit.IsLastNode(circuit.GetLastNode()) {
		t.Logf("Node selected is not the last node")
		t.Fail()
	}
}

//Tests that len returns the correct length
func TestCircuit_Len(t *testing.T) {
	for i := 1; i < 100; i++ {
		circuit := Circuit{nodes: make([]*id.ID, i)}

		if circuit.Len() != i {
			t.Errorf("Circuit.Len: Incorrect length returned,"+
				"Expected: %v, Received: %v", i, circuit.Len())
		}
	}
}

//Tests that all nodes in the circuit return the correct next node
func TestCircuit_GetNextNode(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(9, t)

	circuit := NewCircuit(nodeIdList)

	for index, nid := range nodeIdList {

		expectedNid := nodeIdList[(index+1)%len(nodeIdList)]

		next := circuit.GetNextNode(nid)

		if !expectedNid.Cmp(next) {
			t.Errorf("Circuit.GetNextNode: Returned the incorrect node from index %v,"+
				"Expected: %s, Received: %s", index, expectedNid, next)
		}
	}
}

//Tests GetNextNode panics when the passed node is invalid
func TestCircuit_GetNextNode_Invalid(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(9, t)

	circuit := NewCircuit(nodeIdList)

	invalidNodeID := makeNodeId(77, t)

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	circuit.GetNextNode(invalidNodeID)

	t.Errorf("Circuit.GetNextNode: did not panic with invalid nodeID")
}

//Tests that all nodes in the circuit return the correct next node
func TestCircuit_GetPrevNode(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(9, t)

	circuit := NewCircuit(nodeIdList)

	for index, nid := range nodeIdList {

		var prevLoc int
		if index == 0 {
			prevLoc = len(nodeIdList) - 1
		} else {
			prevLoc = index - 1
		}

		expectedNid := nodeIdList[prevLoc]

		next := circuit.GetPrevNode(nid)

		if !expectedNid.Cmp(next) {
			t.Errorf("Circuit.GetPrevNode: Returned the incorrect node from index %v,"+
				"Expected: %s, Received: %s", index, expectedNid, next)
		}
	}
}

//Tests GetPrevNode panics when the passed node is invalid
func TestCircuit_GetPrevNode_Invalid(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(9, t)

	circuit := NewCircuit(nodeIdList)

	invalidNodeID := makeNodeId(77, t)

	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	circuit.GetPrevNode(invalidNodeID)

	t.Errorf("Circuit.GetPrevNode: did not panic with invalid nodeID")
}

//Test that IsFirstNode is only true when passed first node
func TestCircuit_IsFirstNode(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(23, t)

	circuit := NewCircuit(nodeIdList)

	if !circuit.IsFirstNode(nodeIdList[0]) {
		t.Errorf("Circuit.IsFirstNode: Returned that node at index" +
			"0 is not first node")
	}

	for index, nid := range nodeIdList[1:] {
		if circuit.IsFirstNode(nid) {
			t.Errorf("Circuit.IsFirstNode: Returned that node at index"+
				"%v is first node when it is not", index)
		}
	}
}

//Tests that IsLastNode is only true when passed last node
func TestCircuit_IsLastNode(t *testing.T) {
	nodeIdList := makeTestingNodeIdList(23, t)

	circuit := NewCircuit(nodeIdList)

	if !circuit.IsLastNode(nodeIdList[len(nodeIdList)-1]) {
		t.Errorf("Circuit.IsFirstNode: Returned that node at index"+
			"%v of %v is not last node", len(nodeIdList)-1, len(nodeIdList)-1)
	}

	for index, nid := range nodeIdList[:len(nodeIdList)-2] {
		if circuit.IsLastNode(nid) {
			t.Errorf("Circuit.IsLastNode: Returned that node at index"+
				"%v of %v is last node when it is not", index, len(nodeIdList)-1)
		}
	}
}

// Tests GetOrdering() by checking the position of each rotated node list.
func TestCircuit_GetOrdering(t *testing.T) {
	length := 5
	list := makeTestingNodeIdList(length, t)
	c := NewCircuit(list)
	cs := c.GetOrdering()

	checkShift(t, list, cs[0].nodes, 0)
	checkShift(t, list, cs[1].nodes, 1)
	checkShift(t, list, cs[2].nodes, 2)
	checkShift(t, list, cs[3].nodes, 3)
	checkShift(t, list, cs[4].nodes, 4)
}

// Tests ShiftLeft() by creating list of node IDs, shifting them, and checking
// their position.
func TestShiftLeft(t *testing.T) {
	rotations := 0
	list := makeTestingNodeIdList(5, t)
	newList := shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 1
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 2
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 3
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 4
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 5
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 6
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 10
	list = makeTestingNodeIdList(5, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 1
	list = makeTestingNodeIdList(1, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)

	rotations = 0
	list = makeTestingNodeIdList(1, t)
	newList = shiftLeft(list, rotations)
	checkShift(t, list, newList, rotations)
}

// checkShift checks that each element in list a is in the correct shifted
// position in list b.
func checkShift(t *testing.T, a, b []*id.ID, rotations int) {
	length := len(a)
	for i := 0; i < length; i++ {
		newIndex := (i - (rotations % length) + length) % length
		if !reflect.DeepEqual(a[i], b[newIndex]) {
			t.Errorf("RotateLeft() did not properly shift item #%d to position #%d"+
				"\n\texpected: %#v\n\treceived: %#v",
				i, newIndex, a[i], b[newIndex])
		}
	}
}

//Utility function
func makeTestingNodeIdList(len int, t *testing.T) []*id.ID {
	var nodeIdList []*id.ID

	//build a set of nodeIDs for testing
	for i := 0; i < len; i++ {
		nodeIdBytes := make([]byte, id.ArrIDLen)
		nodeIdBytes[0] = byte(i + 1)
		newNodeId := id.NewIdFromBytes(nodeIdBytes, t)
		nodeIdList = append(nodeIdList, newNodeId)
	}

	return nodeIdList
}

func makeNodeId(b byte, t *testing.T) *id.ID {
	invalidNodeIdBytes := make([]byte, id.ArrIDLen)
	invalidNodeIdBytes[0] = b
	invalidNodeID := id.NewIdFromBytes(invalidNodeIdBytes, t)
	return invalidNodeID
}
