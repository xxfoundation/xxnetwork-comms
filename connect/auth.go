////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles authentication logic for the top-level comms object

package connect

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect/token"
	pb "gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/metadata"
)

// Auth represents an authorization state for a message or host
type Auth struct {
	// Indicates whether authentication was successful
	IsAuthenticated bool
	// The information about the Host that sent the authenticated communication
	Sender *Host
	// reason it isn't authenticated if authentication fails
	Reason string
	// The IP Address (excluding port) for the sending host
	IpAddress string
}

// Perform the client handshake to establish reverse-authentication
// no lock is taken because this is assumed to be done exclusively under the
// send lock taken in ProtoComms.transmit()
func (c *ProtoComms) clientHandshake(host *Host) (err error) {
	ctx, cancel := host.GetMessagingContext()
	defer cancel()
	var result *pb.AssignToken
	if host.connection.IsWeb() {
		wc := host.connection.GetWebConn()
		err = wc.Invoke(ctx, "/messages.Generic/RequestToken", &pb.Ping{}, result)
		if err != nil {
			return err
		}
	} else {
		client := pb.NewGenericClient(host.connection.GetGrpcConn())

		// Send the token request message
		result, err = client.RequestToken(ctx,
			&pb.Ping{})
		if err != nil {
			return errors.New(err.Error())
		}
	}

	remoteToken, err := token.Unmarshal(result.Token)
	if err != nil {
		return errors.Errorf("Failed to unmarshal token: %s", err)
	}

	// Pack the authenticated message with signature enabled
	msg, err := c.PackAuthenticatedMessage(&pb.AssignToken{
		Token: result.Token,
	}, host, true)
	if err != nil {
		return errors.New(err.Error())
	}

	// Add special client-specific info to the
	// newly generated authenticated message if needed
	if c.pubKeyPem != nil && c.salt != nil {
		msg.Client.Salt = c.salt
		msg.Client.PublicKey = string(c.pubKeyPem)
	}

	// Set up the context
	ctx, cancel = host.GetMessagingContext()
	defer cancel()

	if host.connection.IsWeb() {
		wc := host.connection.GetWebConn()
		err = wc.Invoke(ctx, "/messages.Generic/AuthenticateToken", &pb.Ping{}, result)
		if err != nil {
			return err
		}
	} else {
		client := pb.NewGenericClient(host.connection.GetGrpcConn())

		// Send the authenticate token message
		_, err = client.AuthenticateToken(ctx, msg)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	jww.TRACE.Printf("Negotiated Remote token: %v", remoteToken)
	// Assign the host token
	host.transmissionToken.Set(remoteToken)

	return
}

// Convert any message type into a authenticated message
func (c *ProtoComms) PackAuthenticatedMessage(msg proto.Message, host *Host,
	enableSignature bool) (*pb.AuthenticatedMessage, error) {

	// Marshall the provided message into an Any type
	anyMsg, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Build the authenticated message
	authMsg := &pb.AuthenticatedMessage{
		ID:      c.networkId.Marshal(),
		Token:   host.transmissionToken.GetBytes(),
		Message: anyMsg,
		Client: &pb.ClientID{
			Salt:      make([]byte, 0),
			PublicKey: "",
		},
	}

	// If signature is enabled, sign the message and add to payload
	if enableSignature && !c.disableAuth {
		authMsg.Signature, err = c.signMessage(msg, host.GetId())
		if err != nil {
			return nil, err
		}
	}

	return authMsg, nil
}

// Add authentication fields to a given context and return it
func (c *ProtoComms) PackAuthenticatedContext(host *Host,
	ctx context.Context) context.Context {

	ctx = metadata.AppendToOutgoingContext(ctx, "ID", c.networkId.String())
	ctx = metadata.AppendToOutgoingContext(ctx, "TOKEN",
		base64.StdEncoding.EncodeToString(host.transmissionToken.GetBytes()))
	return ctx
}

// Returns authentication packed into a context
func UnpackAuthenticatedContext(ctx context.Context) (*pb.AuthenticatedMessage, error) {
	auth := &pb.AuthenticatedMessage{}
	var err error

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	idStr := md.Get("ID")[0]
	auth.ID, err = base64.StdEncoding.DecodeString(idStr)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decode authentication ID")
	}

	tokenStr := md.Get("TOKEN")[0]
	auth.Token, err = base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, errors.WithMessage(err, "could not decode authentication Live")
	}

	return auth, nil
}

// Generates a new token and adds it to internal state
func (c *ProtoComms) GenerateToken() ([]byte, error) {
	return c.tokens.Generate().Marshal(), nil
}

// Validates a signed token using internal state
func (c *ProtoComms) ValidateToken(msg *pb.AuthenticatedMessage) (err error) {
	// Convert EntityID to ID
	senderId, err := id.Unmarshal(msg.ID)
	if err != nil {
		return err
	}

	// Verify the Host exists for the provided ID
	host, ok := c.GetHost(senderId)
	if !ok {
		return errors.Errorf("No host set up with %s, refusing contact", senderId)
	}

	// Get the signed token
	tokenMsg := &pb.AssignToken{}
	err = ptypes.UnmarshalAny(msg.Message, tokenMsg)
	if err != nil {
		return errors.Errorf("Unable to unmarshal token: %+v", err)
	}

	remoteToken, err := token.Unmarshal(tokenMsg.Token)
	if err != nil {
		return errors.Errorf("Unable to unmarshal token: %+v", err)
	}

	// Verify the token signature unless disableAuth has been set for testing
	if !c.disableAuth {
		if err := c.verifyMessage(tokenMsg, msg.Signature, host); err != nil {
			return errors.Errorf("Invalid token signature: %+v", err)
		}
	}

	ok = c.tokens.Validate(remoteToken)
	if !ok {
		jww.ERROR.Printf("Failed to validate token %v from %s", remoteToken, host)
		return errors.Errorf("Failed to validate token: %v", remoteToken)
	}
	// Token has been validated and can be safely stored
	host.receptionToken.Set(remoteToken)
	jww.DEBUG.Printf("Live validated: %v", tokenMsg.Token)
	return
}

// AuthenticatedReceiver handles reception of an AuthenticatedMessage,
// checking if the host is authenticated & returning an Auth state
func (c *ProtoComms) AuthenticatedReceiver(msg *pb.AuthenticatedMessage, ctx context.Context) (*Auth, error) {

	// Retrieve the IP address
	ipAddr, _, err := GetAddressFromContext(ctx)
	if err != nil {
		return &Auth{}, errors.Errorf("Failed to retrieve ip address from context: %v", err)
	}

	// Convert EntityID to ID
	msgID, err := id.Unmarshal(msg.ID)
	if err != nil {
		return &Auth{
			IsAuthenticated: false,
			Sender:          &Host{},
			IpAddress:       ipAddr,
			Reason: fmt.Sprintf("Host {%v} cannot be "+
				"unmarshaled: %s", msg.ID, err),
		}, nil
	}

	// Try to obtain the Host for the specified ID
	host, ok := c.GetHost(msgID)
	if !ok {
		return &Auth{
			IsAuthenticated: false,
			Sender:          &Host{},
			IpAddress:       ipAddr,
			Reason:          fmt.Sprintf("Host {%s} cannot be found", msgID),
		}, nil
	}

	remoteToken, err := token.Unmarshal(msg.Token)
	if err != nil {
		return &Auth{
			IsAuthenticated: false,
			Sender:          host,
			IpAddress:       ipAddr,
			Reason:          fmt.Sprintf("Token {%v} cannot be unmarshaled", msg.Token),
		}, nil
	}

	// get the hosts reception token
	receptionToken, ok := host.receptionToken.Get()
	if !ok {
		return &Auth{
			IsAuthenticated: false,
			Sender:          host,
			IpAddress:       ipAddr,
			Reason: fmt.Sprintf("failed to authenticate token %v, "+
				"no reception token for %s", remoteToken, host.id),
		}, nil
	}

	// check if the tokens are the same
	if !receptionToken.Equals(remoteToken) {
		return &Auth{
			IsAuthenticated: false,
			Sender:          host,
			IpAddress:       ipAddr,
			Reason: fmt.Sprintf("failed to authenticate token %v, "+
				"does not match reception token %v for %s", remoteToken,
				receptionToken, host.id),
		}, nil
	}

	// Assemble the Auth object
	res := &Auth{
		IsAuthenticated: true,
		Sender:          host,
		IpAddress:       ipAddr,
		Reason:          "authenticated",
	}

	jww.TRACE.Printf("Authentication status: %v, ProvidedId: %v ProvidedToken: %v",
		res.IsAuthenticated, msg.ID, msg.Token)
	return res, nil
}

// DisableAuth makes the authentication code skip signing and signature verification if the
// set.  Can only be set while in a testing structure.  Is not thread safe.
func (c *ProtoComms) DisableAuth() {
	jww.WARN.Print("Auth checking disabled, running insecurely")
	c.disableAuth = true
}

// Takes a message and returns its signature
// The message is signed with the ProtoComms RSA PrivateKey
func (c *ProtoComms) signMessage(msg proto.Message, recipientID *id.ID) ([]byte, error) {
	// Hash the message data
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	hash.Write(msgBytes)
	// Hash in the ID of the intended recipient. This prevents potential
	// replay attacks
	hash.Write(recipientID.Bytes())
	hashed := hash.Sum(nil)

	jww.TRACE.Printf("SignMessage: Signing for host ID [%v]", recipientID)
	jww.TRACE.Printf("SignMessage: hash data: %v", hashed)
	jww.TRACE.Printf("SignMessage: Hashed with ID: %v", recipientID)

	// Obtain the private key
	key := c.GetPrivateKey()
	if key == nil {
		return nil, errors.Errorf("Cannot sign message: No private key")
	}
	// Sign the message and return the signature
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return signature, nil
}

// Takes a message and a Host, verifies the signature
// using Host public key, returning an error if invalid
func (c *ProtoComms) verifyMessage(msg proto.Message, signature []byte, host *Host) error {

	// Deal with edge case in which gateways and servers
	// haven't added each other as hosts yet, and dealing with
	// temporary or dummy ID's. Specifically mitigates
	// gateway's requesting the permissioning address for it's node.
	// fixme: there should be a better way of doing this. Suggested solutions:
	//  - Adding the permissioning address to gateway's config file
	//      (difficult/breaking on deployment side)
	//  - Having the gateway generate the Node's ID on startup. The node's
	//    ID could be pregenerated, ie generated along with it's certs
	//    and added to the Node field in gateway's config (also deployment side)
	var idToHash *id.ID
	if host.id.Cmp(&id.DummyUser) || host.id.Cmp(&id.TempGateway) {
		idToHash = id.DummyUser.DeepCopy()
		idToHash.SetType(id.Node)
	} else {
		idToHash = c.networkId
	}

	// Get hashed data of the message
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return errors.New(err.Error())
	}
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	hash.Write(msgBytes)
	// Hash in the ID of the intended recipient. This prevents potential
	// replay attacks
	hash.Write(idToHash.Bytes())
	hashed := hash.Sum(nil)

	jww.TRACE.Printf("VerifyMessage: Verifying host ID [%v]", host.GetId())
	jww.TRACE.Printf("VerifyMessage: hash data: %v", hashed)
	jww.TRACE.Printf("VerifyMessage: Hashed with ID: %v", idToHash)

	// Verify signature of message using host public key
	err = rsa.Verify(host.rsaPublicKey, options.Hash, hashed, signature, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
