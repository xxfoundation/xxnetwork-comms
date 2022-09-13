package connect

import (
	"context"
	"gitlab.com/xx_network/comms/connect/token"
	pb "gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

type TestGenericServer struct {
	resp string
}

func (ts *TestGenericServer) AuthenticateToken(context.Context, *pb.AuthenticatedMessage) (*pb.Ack, error) {
	return &pb.Ack{Error: ts.resp}, nil
}

func (ts *TestGenericServer) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	return &pb.AssignToken{Token: []byte(ts.resp)}, nil
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
		pb.RegisterGenericServer(s, &TestGenericServer{resp: "response"})
		pc := ProtoComms{
			networkId:   id.NewIdFromString("zezima", id.User, t),
			disableAuth: true,
			tokens:      token.NewMap(),
			Manager:     newManager(),
			netListener: lis,
			grpcServer:  s,
		}
		pc.ServeWithWeb()
	}()
	time.Sleep(time.Second * 5)

	err = h.connect()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := h.GetMessagingContext()
	defer cancel()

	resp := &pb.Ack{}
	err = h.connection.GetWebConn().Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.AuthenticatedMessage{}, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Error)
}

// TODO: Re-enable once HTTPS is supported
func TestWebConnection_TLS(t *testing.T) {
	addr := "0.0.0.0:11421"

	certBytes, err := utils.ReadFile(testkeys.GetNodeCertPath())
	if err != nil {
		t.Fatal(err)
	}

	keyBytes, err := utils.ReadFile(testkeys.GetNodeKeyPath())
	if err != nil {
		t.Fatal(err)
	}

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

	h, err := newHost(hostId, addr, certBytes, hostParams)
	if err != nil {
		t.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGenericServer(s, &TestGenericServer{})

	pk, err := rsa.LoadPrivateKeyFromPem(keyBytes)
	if err != nil {
		t.Fatal(err)
	}
	salt := make([]byte, 8)
	_, err = rng.Read(salt)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		pc := ProtoComms{
			networkId:   id.NewIdFromString("zezima", id.User, t),
			privateKey:  pk,
			disableAuth: false,
			tokens:      token.NewMap(),
			Manager:     newManager(),
			netListener: lis,
			grpcServer:  s,
			pubKeyPem:   certBytes,
			salt:        nil,
		}
		pc.ServeWithWeb()
	}()
	time.Sleep(time.Second * 5)

	err = h.connect()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := h.GetMessagingContext()
	defer cancel()

	resp := &pb.Ack{}
	err = h.connection.GetWebConn().Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.AuthenticatedMessage{}, resp)
	if err != nil {
		t.Fatal(err)
	}
}
