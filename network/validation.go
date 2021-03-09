///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package network

// Level of validation types for pulling our round structure
// 	Strict: Signatures are checked every time (intended for nodes)
//	Lazy: Only check we're involved, only verifies the first retrieval
//	None: no signature checks are done
type ValidationType uint8
const(
	Strict ValidationType = iota
	Lazy
	None
)
