////////////////////////////////////////////////////////////////////////////////
// Copyright © 2024 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package testkeys

import _ "embed"

//go:embed cmix.rip.crt
var nodeCert []byte

//go:embed cmix.rip.key
var nodeKey []byte

//go:embed gateway.cmix.rip.crt
var gatewayCert []byte

//go:embed gateway.cmix.rip.key
var gatewayKey []byte

// These functions are used to cover TLS connection code in tests.

// GetNodeCert returns the contents of cmix.rip.crt.
func GetNodeCert() []byte { return nodeCert }

// GetNodeKey returns the contents of cmix.rip.key.
func GetNodeKey() []byte { return nodeKey }

// GetGatewayCert returns the contents of gateway.cmix.rip.crt.
func GetGatewayCert() []byte { return gatewayCert }

// GetGatewayKey returns the contents of gateway.cmix.rip.key.
func GetGatewayKey() []byte { return gatewayKey }
