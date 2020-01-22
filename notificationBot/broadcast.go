////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package notificationBot

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// NotificationBot -> Permissioning
func (nb *Comms) RequestNdf(host *connect.Host, message *pb.NDFHash) (*pb.NDF, error) {

	// Call the ProtoComms RequestNdf call
	return nb.ProtoComms.RequestNdf(host, message)
}

/* TODO: Uncomment when XX-1777 is done (Add PollForNotifications to Gateways)
// Notification Bot -> Gateway
func (nb *Comms) PollNotifications(host *connect.Host, message *pb.Ping) (*pb.IDList, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		authMsg, err := nb.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Send the message
		resultMsg, err := pb.GatewayClient(
			conn).PollNotifications(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending Request Notification message: %+v", message)
	resultMsg, err := nb.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.IDList{}
	return result, ptypes.UnmarshalAny(resultMsg, result)


}
*/
