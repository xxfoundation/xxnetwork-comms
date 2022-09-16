////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/xx_network/primitives/id"
	"sync"
)

//structure which holds a list of IP address to override
type IpOverrideList struct {
	ipOverride map[id.ID]string
	sync.Mutex
}

//creates a new list over IP overrides
func NewIpOverrideList() *IpOverrideList {
	return &IpOverrideList{
		ipOverride: make(map[id.ID]string),
	}
}

//sets an id to be overridden with a specific IP address
func (iol *IpOverrideList) Override(oid *id.ID, ip string) {
	iol.Lock()
	iol.ipOverride[*oid] = ip
	iol.Unlock()
}

// checks if an ip should be overwritten. returns the passed IP if it should not
// be overwritten
func (iol *IpOverrideList) CheckOverride(cid *id.ID, ip string) string {
	iol.Lock()
	defer iol.Unlock()
	if oip, exists := iol.ipOverride[*cid]; exists {
		return oip
	}
	return ip
}
