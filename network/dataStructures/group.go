////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

// Define the network.Instance's group object and methods for creation and updating

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/ndf"
	"sync"
	"testing"
)

// Struct that handles and updates cyclic.Groups
type Group struct {
	groupString string
	cyclicGroup *cyclic.Group
	*sync.RWMutex
}

// NewGroup creates a ds.Group with a cyclic.Group and a mutex
func NewGroup() *Group {
	return &Group{
		RWMutex: &sync.RWMutex{},
	}
}

// Get returns the ds.Groups's cyclic group
func (g *Group) Get() *cyclic.Group {
	return g.cyclicGroup
}

// Get returns the ds.Groups's cyclic group string
func (g *Group) GetString() string {
	return g.groupString
}

// Update sets the group's string and cyclic.Group object
// If these values have not been set yet, we set these two values
// If these values are set and the newGroup is different, it errors
//  as the group should be immutable after being set
func (g *Group) Update(newGroup string) error {
	g.Lock()
	defer g.Unlock()
	// Check if groupString has not been set
	if g.groupString == "" {
		// If so initialize these values
		g.groupString = newGroup

		// Create cyclic.Group object
		grp, err := toGroup(newGroup)
		if err != nil {
			return errors.Errorf("Unable to update group: %+v", err)
		}
		// Set value
		g.cyclicGroup = grp
	} else if g.groupString != newGroup {
		// If they have already been set and the newGroup is a different value,
		return errors.Errorf("Attempt to modify an already initialized group")
	}

	return nil
}

// Utility function for NewInstanceTesting that directly sets cyclic.Group object
// USED FOR TESTING PURPOSED ONLY
func (g Group) UpdateCyclicGroupTesting(group *cyclic.Group, i interface{}) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		panic("Should not be able to directly set cyclic group outside of testing purposes")
	}

	g.cyclicGroup = group
}

// toGroup is a helper function which converts a string representing a cyclic.Group
// into a cyclic.Group object
func toGroup(gprString string) (*cyclic.Group, error) {
	// Convert it into an ndf group
	tmpNdf := &ndf.Group{}
	err := json.Unmarshal([]byte(gprString), tmpNdf)
	if err != nil {
		return nil, errors.Errorf("Unable to marshal new group: %+v", err)
	}

	// Pull out the prime and generator values from the ndf.Group object
	prime := large.NewIntFromString(tmpNdf.Prime, 16)
	generator := large.NewIntFromString(tmpNdf.Generator, 16)

	// Create the group with the above values
	return cyclic.NewGroup(prime, generator), nil
}
