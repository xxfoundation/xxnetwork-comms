////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"testing"
	"os"
	"math/rand"
	"fmt"
)

var ServerAddress = ""

func TestMain(m *testing.M) {
	ServerAddress = fmt.Sprintf("localhost:%d", (rand.Intn(2000) + 4000))
	os.Exit(m.Run())
}
