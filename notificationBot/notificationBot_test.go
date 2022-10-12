////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package notificationBot

import (
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelTrace)
	connect.TestingOnlyDisableTLS = true
	os.Exit(m.Run())
}

var botPortLock sync.Mutex
var botPort = 1500

// Helper function to prevent port collisions
func getNextAddress() string {
	botPortLock.Lock()
	defer func() {
		botPort++
		botPortLock.Unlock()
	}()
	return fmt.Sprintf("0.0.0.0:%d", botPort)
}

// Error path: Start bot with bad certs
func TestBadCerts(t *testing.T) {
	testID := id.NewIdFromString("test", id.Generic, t)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Address := getNextAddress()

	// This should panic and cause the defer func above to run
	_ = StartNotificationBot(testID, Address, NewImplementation(),
		[]byte("bad cert"), []byte("bad key"))
}
