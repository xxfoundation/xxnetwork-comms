package message

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"

	jww "github.com/spf13/jwalterweatherman"
)

// Send an AskOnline message to a particular server
func SendAskOnline(addr string, message *pb.Ping) (*pb.Pong, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		jww.ERROR.Printf("Failed to connect to server with address %s",
			addr)
	}

	c := pb.NewMixMessageServiceClient(conn)
	// Send AskOnline Request and check that we get an AskOnlineAck back
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	result, err := c.AskOnline(ctx, &pb.Ping{})
	if err != nil {
		jww.ERROR.Printf("AskOnline: Error received: %s", err)
	} else {
		jww.INFO.Printf("AskOnline: %v is online!", addr)
	}
	cancel()
	conn.Close()

	return result, nil
}
