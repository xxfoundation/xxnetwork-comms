package mixmessages

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"testing"
)

// Unit test for helper function to check connectivity of a given address
// To run test on osx, in terminal run:
// nc -l 8080
//func TestCheckConn(t *testing.T) {
//	ok := checkConn("www.google.com:80")
//	if !ok {
//		t.Error("Failed to connect")
//	}
//}

// Mock the Addr interface so we can set the address to check for testing
type MockAddr struct{}

// This does not get used, just fills the interface
func (m *MockAddr) Network() string {
	return ""
}

// Hit an open port on google
func (m *MockAddr) String() string {
	return "www.google.com:80"
}

// Test the functionality & returns of CheckConnectivity
func TestCheckConnectivity(t *testing.T) {
	ctx := context.Background()
	p := &peer.Peer{
		Addr: &MockAddr{},
	}
	ctx = peer.NewContext(ctx, p)
	resp, err := CheckConnectivity(ctx, "", 0)
	if err != nil {
		t.Errorf("Failed to check connectivity with default params: %+v", err)
	}
	if resp.CallerAvailable != true && resp.OtherAvailable != false {
		t.Errorf("Did not get expected response for default params: %+v", resp)
	}

	resp, err = CheckConnectivity(ctx, "", 80)
	if err != nil {
		t.Errorf("Failed to check connectivity with optional port: %+v", err)
	}
	if resp.CallerAvailable != true && resp.OtherAvailable != false {
		t.Errorf("Did not get expected response for optional port: %+v", resp)
	}

	resp, err = CheckConnectivity(ctx, "www.google.com", 80)
	if err != nil {
		t.Errorf("Failed to check connectivity with optional address: %+v", err)
	}
	if resp.CallerAvailable != true && resp.OtherAvailable != true {
		t.Errorf("Did not get expected response for optional address: %+v", resp)
	}
}
