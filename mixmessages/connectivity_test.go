package mixmessages

import (
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

// Test the functionality & returns of CheckConnectivity
func TestCheckConnectivity(t *testing.T) {
	resp, err := CheckConnectivity("www.google.com", "80", "", "0")
	if err != nil {
		t.Errorf("Failed to check connectivity with default params: %+v", err)
	}
	if resp.CallerAvailable != true && resp.OtherAvailable != false {
		t.Errorf("Did not get expected response for default params: %+v", resp)
	}

	resp, err = CheckConnectivity("www.google.com", "80", "", "80")
	if err != nil {
		t.Errorf("Failed to check connectivity with optional port: %+v", err)
	}
	if resp.CallerAvailable != true && resp.OtherAvailable != false {
		t.Errorf("Did not get expected response for optional port: %+v", resp)
	}

	resp, err = CheckConnectivity("www.google.com", "80", "www.google.com", "80")
	if err != nil {
		t.Errorf("Failed to check connectivity with optional address: %+v", err)
	}
	if resp.CallerAvailable != true && resp.OtherAvailable != true {
		t.Errorf("Did not get expected response for optional address: %+v", resp)
	}
}
