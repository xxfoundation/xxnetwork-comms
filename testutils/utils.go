///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package testutils

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"testing"
)

func LoadKeyTesting(t *testing.T) *rsa.PublicKey {
	if t == nil {
		jww.FATAL.Panicf("LoadKeyTesting is a testing only function")
	}

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)

	privKey, err := rsa.LoadPrivateKeyFromPem(keyData)
	if err != nil {
		t.Errorf("Could not load public key: %v", err)
		t.FailNow()
	}

	return privKey.GetPublic()
}

// Utility function which signs a round info message
func SignRoundInfo(ri *pb.RoundInfo) error {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)

	privKey, err := rsa.LoadPrivateKeyFromPem(keyData)
	if err != nil {
		return errors.Errorf("Could not load public key: %v", err)
	}


	err = signature.Sign(ri, privKey)
	if err != nil {
		return errors.Errorf("Could not sign round info: %+v", err)
	}
	return nil
}
