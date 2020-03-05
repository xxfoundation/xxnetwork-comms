////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package dataStructures

import (
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/large"
	"sync"
)

// todo docstring
type Group struct {
	grp   *cyclic.Group
	mutex *sync.RWMutex
}

// NewGroup creates a ds.Group with a cyclic.Group and a mutex
func NewGroup(p, g *large.Int) *Group {
	ourGroup := cyclic.NewGroup(p, g)

	return &Group{
		grp:   ourGroup,
		mutex: &sync.RWMutex{},
	}
}

// Get returns the ds.Groups's cyclic group
func (g *Group) Get() *cyclic.Group {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.grp
}
