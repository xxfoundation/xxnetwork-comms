package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/xx_network/comms/connect"
	"git.xx.network/xx_network/comms/messages"
	"google.golang.org/grpc"
)

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendDeleteMessage(host *connect.Host, message *pb.FactRemovalRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Pack our authenticated message
		authMsg, err := c.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RemoveFact(ctx, authMsg)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery Register User Function
func (c *Comms) SendRegisterUser(host *connect.Host, message *pb.UDBUserRegistration) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Pack our authenticated message
		authMsg, err := c.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RegisterUser(ctx, authMsg)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery Register Fact Function
func (c *Comms) SendRegisterFact(host *connect.Host, message *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Pack our authenticated message
		authMsg, err := c.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RegisterFact(ctx, authMsg)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Register Fact message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.FactRegisterResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendConfirmFact(host *connect.Host, message *pb.FactConfirmRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Pack our authenticated message
		authMsg, err := c.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).ConfirmFact(ctx, authMsg)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}
