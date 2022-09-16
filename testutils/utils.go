////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package testutils

import (
	"context"
	"crypto/rand"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/ec"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"google.golang.org/grpc/peer"
	"net"
	"testing"
	"time"
)

const privKeyEncoded = `uVAt6d+y3XW699L3THlcoTA2utw2dhoqnX6821x6OcnOliwX84eajmp45IZ+STw0dUl8uJtZwDKDuHVX6ZpGzg==`

func LoadPublicKeyTesting(i interface{}) (*rsa.PublicKey, error) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	privKey, err := LoadPrivateKeyTesting(i)
	if err != nil {
		return nil, errors.Errorf("Could not load private key: %v", err)
	}

	return privKey.GetPublic(), nil
}

func LoadPrivateKeyTesting(i interface{}) (*rsa.PrivateKey, error) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)

	privKey, err := rsa.LoadPrivateKeyFromPem(keyData)
	if err != nil {
		return nil, errors.Errorf("Could not load public key: %v", err)
	}

	return privKey, nil

}

func LoadEllipticPublicKey(i interface{}) (*ec.PrivateKey, error) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	ecKey, err := ec.NewKeyPair(rand.Reader)
	if err != nil {
		return nil, errors.Errorf("Failed to generate new keypair: %v", err)
	}
	err = ecKey.UnmarshalText(privKeyEncoded)
	if err != nil {
		return nil, errors.Errorf("Failed to unmarshal private key: %v", err)
	}
	return ecKey, nil

}

// Utility function which signs a round info message
func SignRoundInfoRsa(ri *pb.RoundInfo, i interface{}) error {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoRsa is restricted to testing only. Got %T", i)
	}

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)

	privKey, err := rsa.LoadPrivateKeyFromPem(keyData)
	if err != nil {
		return errors.Errorf("Could not load public key: %v", err)
	}

	err = signature.SignRsa(ri, privKey)
	if err != nil {
		return errors.Errorf("Could not sign round info: %+v", err)
	}
	return nil
}

func SignRoundInfoEddsa(ri *pb.RoundInfo, key *ec.PrivateKey, i interface{}) error {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoEddsa is restricted to testing only. Got %T", i)
	}
	err := signature.SignEddsa(ri, key)
	if err != nil {
		return errors.Errorf("Could not sign round info: %+v", err)
	}
	return nil

}

// NewContextTesting constructs a context.Context object on
// the local Unix default domain (UDP) port
func NewContextTesting(i interface{}) (context.Context, context.CancelFunc) {
	switch i.(type) {
	case *testing.T:
		break
	case *testing.M:
		break
	case *testing.B:
		break
	default:
		jww.FATAL.Panicf("SignRoundInfoEddsa is restricted to testing only. Got %T", i)
	}

	protoCtx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("udp", "0.0.0.0:53", timeout)
	if err != nil {
		jww.FATAL.Fatalf("Failed to get a conn object in setup: %v", err)
	}
	p := &peer.Peer{
		Addr: conn.RemoteAddr(),
	}

	return peer.NewContext(protoCtx, p), cancel
}
