package message

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"

	jww "github.com/spf13/jwalterweatherman"
)

// Send a NetworkError message to a particular server
func SendNetworkError(addr string, message *pb.ErrorMessage) (*pb.ErrorAck, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		jww.ERROR.Printf("Failed to connect to server with address %s",
			addr)
	}

	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	result, err := c.NetworkError(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("NetworkError: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}
