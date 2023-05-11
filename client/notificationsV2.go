package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// RegisterTrackedID registers the given ID to be tracked. The request is signed
// Returns an error if TransmissionRSA is not registered with a valid token.
// The actual ID is not revealed, instead an intermediary value is sent which cannot
// be revered to get the ID, but is repeatable. So it can be rainbow-tabled.
func (c *Comms) RegisterTrackedID(host *connect.Host, message *pb.RegisterTrackedIdRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var err error
		var resultMsg = &messages.Ack{}
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.NotificationBot/RegisterTrackedID",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewNotificationBotClient(conn.GetGrpcConn()).
				RegisterTrackedID(ctx, message)
		}
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending RegisterTrackedID message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// UnregisterTrackedID unregisters the given tracked ID. The request is signed.
// Does not return an error if the token cannot be found
func (c *Comms) UnregisterTrackedID(host *connect.Host, message *pb.UnregisterTrackedIdRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var err error
		var resultMsg = &messages.Ack{}
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.NotificationBot/UnregisterTrackedID",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewNotificationBotClient(conn.GetGrpcConn()).
				UnregisterTrackedID(ctx, message)
		}
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending UnregisterTrackedID message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// RegisterToken registers the given token. It evaluates that the TransmissionRsaRegistarSig is
// correct. The RSA->PEM relationship is one to many. It will succeed if the token is already
// registered.
func (c *Comms) RegisterToken(host *connect.Host, message *pb.RegisterTokenRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var err error
		var resultMsg = &messages.Ack{}
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.NotificationBot/RegisterToken",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewNotificationBotClient(conn.GetGrpcConn()).
				RegisterToken(ctx, message)
		}
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending RegisterToken message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// UnregisterToken unregisters the given token. The request is signed.
// Does not return an error if the token cannot be found
func (c *Comms) UnregisterToken(host *connect.Host, message *pb.UnregisterTokenRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		var err error
		var resultMsg = &messages.Ack{}
		if conn.IsWeb() {
			wc := conn.GetWebConn()
			err = wc.Invoke(ctx, "/mixmessages.NotificationBot/UnregisterToken",
				message, resultMsg)
		} else {
			resultMsg, err = pb.NewNotificationBotClient(conn.GetGrpcConn()).
				UnregisterToken(ctx, message)
		}
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending UnregisterToken message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
