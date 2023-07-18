package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// Login to the server, receiving an authentication token
func (rc *Comms) Login(host *connect.Host, msg *pb.RsAuthenticationRequest) (*pb.RsAuthenticationResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Send the message
		resultMsg, err := pb.NewRemoteSyncClient(conn.GetGrpcConn()).
			Login(ctx, msg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := rc.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RsAuthenticationResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Read a resource from a RemoteSync server.
func (rc *Comms) Read(host *connect.Host, msg *pb.RsReadRequest) (*pb.RsReadResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Send the message
		resultMsg, err := pb.NewRemoteSyncClient(conn.GetGrpcConn()).
			Read(ctx, msg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := rc.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RsReadResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Write data to a path at a RemoteSync server
func (rc *Comms) Write(host *connect.Host, msg *pb.RsWriteRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Send the message
		resultMsg, err := pb.NewRemoteSyncClient(conn.GetGrpcConn()).
			Write(ctx, msg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := rc.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &messages.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// GetLastModified returns the last time a path was modified.
func (rc *Comms) GetLastModified(host *connect.Host, msg *pb.RsReadRequest) (*pb.RsTimestampResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Send the message
		resultMsg, err := pb.NewRemoteSyncClient(conn.GetGrpcConn()).
			GetLastModified(ctx, msg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := rc.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RsTimestampResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// GetLastWrite returns the last time a remote sync server was modified.
func (rc *Comms) GetLastWrite(host *connect.Host, msg *pb.RsLastWriteRequest) (*pb.RsTimestampResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Send the message
		resultMsg, err := pb.NewRemoteSyncClient(conn.GetGrpcConn()).
			GetLastWrite(ctx, msg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := rc.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RsTimestampResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// ReadDir returns all entries in a given path.
func (rc *Comms) ReadDir(host *connect.Host, msg *pb.RsReadRequest) (*pb.RsReadDirResponse, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		// Send the message
		resultMsg, err := pb.NewRemoteSyncClient(conn.GetGrpcConn()).
			ReadDir(ctx, msg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := rc.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RsReadDirResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
