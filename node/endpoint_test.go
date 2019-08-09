package node

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"gitlab.com/elixxir/comms/testkeys"
	"testing"
)

func TestDownloadTopology(t *testing.T) {
	ctx := context.Background()

	wrapped := mixmessages.NodeTopology{
		Topology: []*mixmessages.NodeInfo{},
	}

	a, _ := ptypes.MarshalAny(&wrapped)

	msg := mixmessages.SignedMessage{
		ID:        "Permissioning",
		Signature: []byte("test"),
		Message:   a,
	}

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	ServerAddress := getNextServerAddress()
	RegAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(),
		certData, keyData)
	reg := registration.StartRegistrationServer(RegAddress,
		registration.NewImplementation(), certData, keyData)
	defer server.Shutdown()
	defer reg.Shutdown()

	regID := MockID("Permissioning")
	err := server.ConnectToRegistration(regID, RegAddress, certData)

	if err != nil {
		t.Errorf("Download topology could not connect to registration: %+v", err)
	}

	_, err = server.DownloadTopology(ctx, &msg)
	if err != nil {
		t.Errorf("Download topology failed: %+v", err)
	}

}
