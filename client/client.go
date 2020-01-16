package client

import "gitlab.com/elixxir/comms/connect"

// Client object used to implement endpoints and top-level comms functionality
type Comms struct {
	connect.ProtoComms
}
