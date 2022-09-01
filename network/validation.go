////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package network

import "fmt"

// Level of validation types for pulling our round structure
// 	Strict: Signatures are checked every time (intended for nodes)
//	Lazy: Only check we're involved, only verifies the first retrieval
//	None: no signature checks are done
type ValidationType uint8

const (
	Strict ValidationType = iota
	Lazy
	None
)

// Stringer for ValidationType
func (s ValidationType) String() string {
	switch s {
	case Strict:
		return "Strict"
	case Lazy:
		return "Lazy"
	case None:
		return "None"
	default:
		return fmt.Sprintf("UNKNOWN STATE: %d", s)
	}
}
