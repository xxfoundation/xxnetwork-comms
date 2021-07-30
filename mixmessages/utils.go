///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains utils functions for comms

package mixmessages

import jww "github.com/spf13/jwalterweatherman"

// Headers for streaming
const PostPhaseHeader = "batchinfo"
const UnmixedBatchHeader = "unmixedbatchinfo"
const MixedBatchHeader = "mixedBatchInfo"

func DebugMode() {
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)
}

func TraceMode() {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)
}
