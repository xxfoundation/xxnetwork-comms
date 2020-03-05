////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package dataStructures

// todo docstring
type Group struct {
	grp string
}

// NewGroup creates a ds.Group with a cyclic.Group and a mutex
func NewGroup(ourGroup string) *Group {
	return &Group{
		grp: ourGroup,
	}
}

// Get returns the ds.Groups's cyclic group
func (g *Group) Get() string {
	return g.grp
}
