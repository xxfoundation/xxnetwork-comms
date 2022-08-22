package connect

import (
	"context"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	pb "gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"testing"
	"time"
)

type TestGenericServer struct {
}

func (ts *TestGenericServer) AuthenticateToken(context.Context, *pb.AuthenticatedMessage) (*pb.Ack, error) {
	return &pb.Ack{}, nil
}

func (ts *TestGenericServer) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	return &pb.AssignToken{Token: []byte("testtoken")}, nil
}

func TestWebConnection(t *testing.T) {
	addr := "0.0.0.0:11420"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	rng := csprng.NewSystemRNG()
	hostId, err := id.NewRandomID(rng, id.User)
	if err != nil {
		t.Fatal(err)
	}
	hostParams := GetDefaultHostParams()
	TestingOnlyDisableTLS = true
	hostParams.ConnectionType = Web

	h, err := newHost(hostId, addr, nil, hostParams)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		s := grpc.NewServer()
		pb.RegisterGenericServer(s, &TestGenericServer{})
		ws := grpcweb.WrapServer(s, grpcweb.WithOriginFunc(func(origin string) bool { return true }))
		if err := http.Serve(lis, ws); err != nil {
			fmt.Println(err)
			t.Errorf("failed to serve: %v", err)
		}
	}()
	time.Sleep(time.Second * 5)

	err = h.connect()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := h.GetMessagingContext()
	defer cancel()

	// TODO: This fails with RequestToken, seemingly because Ping has no actual contents.  Throws an EOF error when attempting to parse response.  Need to look into this in client repo.
	resp := &pb.AssignToken{}
	err = h.connection.GetWebConn().Invoke(ctx, "/messages.Generic/RequestToken", &pb.Ping{}, resp)
	if err != nil {
		t.Fatal(err)
	}
}
