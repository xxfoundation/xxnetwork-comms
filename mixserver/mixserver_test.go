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

// Smoke test the NetworkError endpoint
func TestNetworkError(t *testing.T) {
	addr := "localhost:5555"
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithInsecure(),
		grpc.WithTimeout(time.Second))
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMixMessageServiceClient(conn)

	// Say hello, check that we get the correct response
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	r, err := c.NetworkError(ctx, &pb.ErrorMessage{Message: "Hello, world!"})
	if err != nil {
		t.Errorf("Error received: %s", err)
	}
	if r.MsgLen != 13 {
		t.Errorf("Expected len of %v, got %v", 13, r)
	}
	defer cancel()
	
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

	// Say hello, check that we get the correct response
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		t.Errorf("could not greet: %v", err)
	} else if r.Message != "Hello MixMessageService" {
		t.Errorf("Wrong greeting: %s", r.Message)
	}
	defer cancel()

	time.Sleep(time.Millisecond*600)
	
	// Send it again, this time expect a timeout error.
	ctx2, cancel2 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	r2, err2 := c.SayHello(ctx2, &pb.HelloRequest{Name: "Ehsy"})
	if err2 == nil {
		t.Errorf("Somehow able to greet: %s", r2.Message)
	}
	defer cancel2()
}
