////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains utils functions for comms

package mixmessages

import jww "github.com/spf13/jwalterweatherman"

func DebugMode() {
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)
}
