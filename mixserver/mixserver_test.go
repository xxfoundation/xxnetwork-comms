// Copyright © 2018 Privategrity Corporation
// All rights reserved.

package mixserver

import (
	"os"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func TestMain(m *testing.M) {
	addr := "localhost:5555"
	go StartServer(addr)
	os.Exit(m.Run())
}

func TestStartServer(t *testing.T) {
	addr := "localhost:5555"
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(),
		grpc.WithTimeout(time.Second))
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMixMessageServiceClient(conn)

	// Contact the server and print out its response.
	name := "MixMessageService"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		t.Errorf("could not greet: %v", err)
	}
	t.Errorf("Greeting: %s", r.Message)
}
