package client

import (
	"gitlab.com/elixxir/comms/connect"
)

// There's no factory method for this, but you can just make an empty variable
// of it and use it and it works
type ClientComms struct {
	connect.ConnectionManager
}
