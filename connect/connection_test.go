package connect

import (
	"context"
	"fmt"
	"gitlab.com/xx_network/comms/connect/token"
	pb "gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"google.golang.org/grpc"
	"testing"
)

type TestGenericServer struct {
	resp string
	pb.UnimplementedGenericServer
}

func (ts *TestGenericServer) AuthenticateToken(context.Context, *pb.AuthenticatedMessage) (*pb.Ack, error) {
	return &pb.Ack{Error: ts.resp}, nil
}

func (ts *TestGenericServer) RequestToken(context.Context, *pb.Ping) (*pb.AssignToken, error) {
	return &pb.AssignToken{Token: []byte(ts.resp)}, nil
}

func TestWebConnection(t *testing.T) {
	addr := "0.0.0.0:11420"

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

	grpcHostParams := GetDefaultHostParams()
	grpcHost, err := newHost(hostId, addr, nil, grpcHostParams)

	errCh := make(chan error)
	go func() {
		s := grpc.NewServer()
		pb.RegisterGenericServer(s, &TestGenericServer{resp: "response"})
		pc := ProtoComms{
			networkId:        id.NewIdFromString("zezima", id.User, t),
			disableAuth:      true,
			tokens:           token.NewMap(),
			Manager:          newManager(),
			listeningAddress: addr,
			grpcServer:       s,
		}

		pc.ServeWithWeb()
		errCh <- pc.ServeHttps(nil, nil)
	}()
	err = <-errCh
	if err != nil {
		t.Fatal(err)
	}

	err = h.connect()
	if err != nil {
		t.Fatal(err)
	}

	err = grpcHost.connect()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := grpcHost.GetMessagingContext()
	resp := &pb.Ack{}
	err = grpcHost.connection.GetGrpcConn().Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.AuthenticatedMessage{}, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Error)
	cancel()

	ctx, cancel = h.GetMessagingContext()
	defer cancel()

	resp = &pb.Ack{}
	err = h.connection.GetWebConn().Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.AuthenticatedMessage{}, resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Error)
}

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

	rng := csprng.NewSystemRNG()
	hostId, err := id.NewRandomID(rng, id.User)
	if err != nil {
		t.Fatal(err)
	}

	pk, err := rsa.LoadPrivateKeyFromPem(keyBytes)
	if err != nil {
		t.Fatal(err)
	}
	salt := make([]byte, 8)
	_, err = rng.Read(salt)
	if err != nil {
		t.Fatal(err)
	}

	pc := ProtoComms{
		networkId:        id.NewIdFromString("zezima", id.User, t),
		privateKey:       pk,
		disableAuth:      false,
		tokens:           token.NewMap(),
		Manager:          newManager(),
		listeningAddress: addr,
		pubKeyPem:        certBytes,
		salt:             nil,
	}

	hostParams := GetDefaultHostParams()
	hostParams.ConnectionType = Web
	h, err := newHost(hostId, addr, certBytes, hostParams)
	if err != nil {
		t.Fatal(err)
	}

	grpcHostParams := GetDefaultHostParams()
	grpcHost, err := newHost(hostId, addr, certBytes, grpcHostParams)
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 5; i++ {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			err = pc.Restart()
			if err != nil {
				t.Fatal(err)
			}

			expectedResponse := fmt.Sprintf("hello! %d", i)
			pb.RegisterGenericServer(pc.grpcServer, &TestGenericServer{resp: expectedResponse})

			pc.ServeWithWeb()
			err = pc.ServeHttps(certBytes, keyBytes)
			if err != nil {
				t.Fatal(err)
			}

			err = h.connect()
			if err != nil {
				t.Fatal(err)
			}

			err = grpcHost.connect()
			if err != nil {
				t.Fatal(err)
			}

			ctx, cancel := grpcHost.GetMessagingContext()
			resp := &pb.Ack{}
			err = grpcHost.connection.GetGrpcConn().Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.AuthenticatedMessage{}, resp)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(resp.Error)
			cancel()

			ctx, cancel = h.GetMessagingContext()
			defer cancel()

			resp = &pb.Ack{}
			err = h.connection.GetWebConn().Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.AuthenticatedMessage{}, resp)
			if err != nil {
				t.Fatalf("Failed to invoke authenticate: %+v", err)
			}
			if resp.Error != expectedResponse {
				t.Errorf("Did not receive expected payload")
			}

			pc.Shutdown()
			h.disconnect()
			grpcHost.disconnect()
		})
	}

}
