package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"google.golang.org/grpc"
)

// Server -> Server CheckConnectivity Function
func (s *Comms) SendCheckConnectivityMessage(host *connect.Host,
	message *pb.Address) (*pb.ConnectivityResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewConnectivityCheckerClient(conn).CheckConnectivity(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Nonce message: %+v", message)
	resultMsg, err := s.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.ConnectivityResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
