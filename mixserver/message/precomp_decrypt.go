package message

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"

	jww "github.com/spf13/jwalterweatherman"
)

func SendPrecompDecrypt(addr string, input *pb.PrecompDecryptMessage) (*pb.Ack, error) {

	// Attempt to connect to nextServer
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Send the PrecompDecrypt message using the Decrypt output
	result, err := c.PrecompDecrypt(ctx, input)
	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompDecrypt: Error received: %s", err)
	}

	return result, err
}
