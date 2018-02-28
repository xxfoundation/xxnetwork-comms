package clusterclient

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"

	jww "github.com/spf13/jwalterweatherman"
)

func SendAskOnline(addr string, message *pb.Ping) (*pb.Pong, error) {
	// Attempt to connect to addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %s",
			addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	// Send the message
	result, err := c.AskOnline(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("AskOnline: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}
