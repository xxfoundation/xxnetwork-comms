package client

import (
	"gitlab.com/elixxir/comms/connect"
	"google.golang.org/grpc"
)

type Client struct {
	manager connect.ConnectionManager
	gs      *grpc.Server
}