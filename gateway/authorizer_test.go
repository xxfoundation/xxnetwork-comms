package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/authorizer"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// Happy path.
func TestComms_SendAuthorizerCertRequest(t *testing.T) {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)

	// Set up gateway
	gwAddr := getNextGatewayAddress()
	gwID := id.NewIdFromString("TestGatewayID", id.Gateway, t)
	gateway := StartGateway(gwID, gwAddr, NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gateway.Shutdown()

	// Set up authorizer
	authAddr := getNextServerAddress()
	authID := &id.Authorizer
	impl := authorizer.NewImplementation()
	receiveChan := make(chan *pb.AuthorizerCertRequest)
	impl.Functions.RequestCert = func(notifBatch *pb.AuthorizerCertRequest) (*messages.Ack, error) {
		go func() { receiveChan <- notifBatch }()
		return &messages.Ack{}, nil
	}
	authServer := authorizer.StartAuthorizerServer(authID, authAddr, impl, nil, nil)
	defer authServer.Shutdown()

	// Create manager and add authorizer as host
	manager := connect.NewManagerTesting(t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(authID, authAddr, nil, params)
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Generate message to send
	msg := &pb.AuthorizerCertRequest{
		Timestamp: uint64(54321),
	}

	// Send auth cert request to authorizer
	resp, err := gateway.SendAuthorizerCertRequest(host, msg)
	if err != nil {
		t.Errorf("SendAuthorizerCertRequest() returned an error: %+v", err)
	}
	if resp == nil {
		t.Errorf("SendAuthorizerCertRequest() did not respond with an AuthorizerCert")
	}

	select {
	case result := <-receiveChan:
		if msg.String() != result.String() {
			t.Errorf("Failed to receive the expected Authorizer Cert."+
				"\nexpected: %s\nreceived: %s", msg, result)
		}
	case <-time.NewTimer(50 * time.Millisecond).C:
		t.Error("Timed out while waiting to receive the Authorizer Cert.")
	}
}

// Happy path.
func TestComms_SendEABCredentialRequest(t *testing.T) {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)

	// Set up gateway
	gwAddr := getNextGatewayAddress()
	gwID := id.NewIdFromString("TestGatewayID", id.Gateway, t)
	gateway := StartGateway(gwID, gwAddr, NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gateway.Shutdown()

	// Set up authorizer
	authAddr := getNextServerAddress()
	authID := &id.Authorizer
	impl := authorizer.NewImplementation()
	receiveChan := make(chan *pb.EABCredentialRequest)
	impl.Functions.SendEABCredentialRequest = func(notifBatch *pb.EABCredentialRequest) (*pb.EABCredentialResponse, error) {
		go func() { receiveChan <- notifBatch }()
		return &pb.EABCredentialResponse{}, nil
	}
	authServer := authorizer.StartAuthorizerServer(authID, authAddr, impl, nil, nil)
	defer authServer.Shutdown()

	// Create manager and add authorizer as host
	manager := connect.NewManagerTesting(t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(authID, authAddr, nil, params)
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Generate message to send
	msg := &pb.EABCredentialRequest{}

	// Send auth EABCredential request to authorizer
	resp, err := gateway.SendEABCredentialRequest(host, msg)
	if err != nil {
		t.Errorf("SendEABCredentialRequest() returned an error: %+v", err)
	}
	if resp == nil {
		t.Errorf("SendEABCredentialRequest() did not respond with an EABCredentialResponse")
	}

	select {
	case result := <-receiveChan:
		if msg.String() != result.String() {
			t.Errorf("Failed to receive the expected Authorizer EABCredential response."+
				"\nexpected: %s\nreceived: %s", msg, result)
		}
	case <-time.NewTimer(50 * time.Millisecond).C:
		t.Error("Timed out while waiting to receive the Authorizer EAB credential response.")
	}
}
