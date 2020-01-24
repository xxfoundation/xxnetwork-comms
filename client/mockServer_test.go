////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains a dummy/mock server instance for testing purposes

package client

import (
	"fmt"
	"sync"
)

var portLock sync.Mutex
var port = 5800

func getNextAddress() string {
	portLock.Lock()
	defer func() {
		port++
		portLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", port)
}
