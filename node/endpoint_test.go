package node

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/elixxir/comms/connect"
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

	testId := "Permissioning"
	msg := mixmessages.SignedMessage{
		ID:        testId,
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
	var manager connect.Manager

	_, err := manager.AddHost(testId, RegAddress, certData, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	if err != nil {
		t.Errorf("Download topology failed: %+v", err)
	}

	_, err = server.DownloadTopology(ctx, &msg)
	if err != nil {
		t.Errorf("Download topology failed: %+v", err)
	}

}
