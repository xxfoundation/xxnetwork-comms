////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var ServerAddress = ""

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	ServerAddress = fmt.Sprintf("localhost:%d", (rand.Intn(2000) + 4000))
	os.Exit(m.Run())
}
