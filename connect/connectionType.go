///////////////////////////////////////////////////////////////////////////////
// Copyright © 2024 xx foundation                                            //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// This file is only compiled for all architectures except WebAssembly.
//go:build !js || !wasm
// +build !js !wasm

package connect

// GetDefaultConnectionType returns Grpc as the default connection type when
// compiling for all architectures except WebAssembly.
func GetDefaultConnectionType() ConnectionType {
	return Grpc
}
