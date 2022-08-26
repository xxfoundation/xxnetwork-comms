package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// Client -> User Discovery Register User Function
func (c *Comms) SendRegisterUser(host *connect.Host, message *pb.UDBUserRegistration) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &messages.Ack{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.UDB/RegisterUser", message, resultMsg)
		} else {
			resultMsg, err = pb.NewUDBClient(conn.GetGrpcConn()).
				RegisterUser(ctx, message)
		}
		if err != nil {
			return nil, err
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
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &pb.FactRegisterResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.UDB/RegisterFact", message, resultMsg)
		} else {
			resultMsg, err = pb.NewUDBClient(conn.GetGrpcConn()).
				RegisterFact(ctx, message)
		}
		if err != nil {
			return nil, err
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
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &messages.Ack{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.UDB/ConfirmFact", message, resultMsg)
		} else {
			resultMsg, err = pb.NewUDBClient(conn.GetGrpcConn()).
				ConfirmFact(ctx, message)
		}
		if err != nil {
			return nil, err
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

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendRemoveFact(host *connect.Host, message *pb.FactRemovalRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &messages.Ack{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.UDB/RemoveFact", message, resultMsg)
		} else {
			resultMsg, err = pb.NewUDBClient(conn.GetGrpcConn()).
				RemoveFact(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete Fact Message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendRemoveUser(host *connect.Host, message *pb.FactRemovalRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &messages.Ack{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.UDB/RemoveUser", message, resultMsg)
		} else {
			resultMsg, err = pb.NewUDBClient(conn.GetGrpcConn()).
				RemoveUser(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete Fact Message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery channel authentication & lease request
func (c *Comms) SendChannelAuthRequest(host *connect.Host, message *pb.ChannelLeaseRequest) (*pb.ChannelLeaseResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var resultMsg = &pb.ChannelLeaseResponse{}
		var err error
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(
				ctx, "/mixmessages.UDB/RequestChannelLease", message, resultMsg)
		} else {
			resultMsg, err = pb.NewUDBClient(conn.GetGrpcConn()).
				RequestChannelLease(ctx, message)
		}
		if err != nil {
			return nil, err
		}
		return ptypes.MarshalAny(resultMsg)
	}

	jww.TRACE.Printf("Sending Channel Auth Request message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.ChannelLeaseResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)

}
