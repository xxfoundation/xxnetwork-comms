////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"io/ioutil"
)

// Send a message to the gateway
func (r *RegistrationComms) SendNodeTopology(id fmt.Stringer,
	message *pb.NodeTopology) error {

	// Attempt to connect to addr
	connection := r.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Wrap message as a generic
	anyMessage, err := ptypes.MarshalAny(message)
	if err != nil {
		jww.ERROR.Printf("ERROR OUT HERE: %+v", err)
	}

	keyBytes, err := ioutil.ReadFile(GlobalKeyPath)
	if err != nil {
		jww.ERROR.Printf("Failed to read private key file at %s: %+v", GlobalKeyPath, err)
	}

	key, err := rsa.LoadPrivateKeyFromPem(keyBytes)
	if err != nil {
		jww.ERROR.Printf("Failed to form private key file from data at %s: %+v", GlobalKeyPath, err)
	}

	// Get hashed data
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	data := []byte(message.String())
	hashed := hash.Sum(data)[len(data):]

	// Sign the thing
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)

	// Form signed message
	signedMessage := pb.SignedMessage{
		Message:   anyMessage,
		Signature: signature,
	}

	// Send the message
	_, err = connection.DownloadTopology(ctx, &signedMessage)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("SendNodeToplogy: Error received: %+v", err)
	}

	cancel()
	return err
}
