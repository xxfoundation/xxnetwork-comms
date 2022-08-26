package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// SendGetPermissioningAddress ping server to return the address of
// permissioning.
func (g *Comms) SendGetPermissioningAddress(host *connect.Host) (string, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn.GetGrpcConn()).
			GetPermissioningAddress(ctx, &messages.Ping{})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending get permissioning address ping.")
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return "", err
	}

	// Marshall the result
	result := &pb.StrAddress{}
	err = ptypes.UnmarshalAny(resultMsg, result)
	if err != nil {
		return "", err
	}
	return result.Address, nil
}
