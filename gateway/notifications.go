package gateway

import (
	"github.com/golang/protobuf/ptypes/any"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
)

// SendNotificationBatch sends the batch of notification data to the
// notification bot.
func (g *Comms) SendNotificationBatch(host *connect.Host, notifBatch *pb.NotificationBatch) error {

	// Create the send function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Pack data into authenticated message
		authMsg, err := g.PackAuthenticatedMessage(notifBatch, host, false)
		if err != nil {
			return nil, err
		}

		// Send the message
		_, err = pb.NewNotificationBotClient(conn.GetGrpcConn()).
			ReceiveNotificationBatch(ctx, authMsg)
		return nil, err
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending notification data batch to notification bot: %s", notifBatch)
	_, err := g.Send(host, f)

	return err
}
