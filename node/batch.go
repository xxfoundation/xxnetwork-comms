///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains logic for batch-related comms

package node

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc/metadata"
)

// UploadUnmixedBatch is the handler for gateway sending a batch to its node
func (s *Comms) UploadUnmixedBatch(server pb.Node_UploadUnmixedBatchServer) error {
	// Extract the authentication info
	authMsg, err := connect.UnpackAuthenticatedContext(server.Context())
	if err != nil {
		return errors.Errorf("Unable to extract authentication info: %+v", err)
	}

	authState, err := s.AuthenticatedReceiver(authMsg, server.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Verify the message authentication
	return s.handler.UploadUnmixedBatch(server, authState)
}

// GetUnmixedBatchStreamHeader gets the header in the metadata from
// the server stream and returns it or an error if it fails.
func GetUnmixedBatchStreamHeader(stream pb.Node_UploadUnmixedBatchServer) (*pb.BatchInfo, error) {

	// Obtain the headers from server metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, errors.New("unable to retrieve meta data / header")
	}

	// Unmarshall the header into a message
	marshledBatch, err := base64.StdEncoding.DecodeString(md.Get(pb.UnmixedBatchHeader)[0])
	if err != nil {
		return nil, err
	}
	batchInfo := &pb.BatchInfo{}
	err = proto.UnmarshalText(string(marshledBatch), batchInfo)
	if err != nil {
		return nil, err
	}

	return batchInfo, nil
}

// ------------------------- DownloadMixedBatch Logic ---------------------------------------- //

// DownloadMixedBatch streams the slots in the completed batch to the gateway
func (s *Comms) DownloadMixedBatch(authMsg *messages.AuthenticatedMessage,
	stream pb.Node_DownloadMixedBatchServer) error {

	authState, err := s.AuthenticatedReceiver(authMsg, stream.Context())
	if err != nil {
		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	//Unmarshall the any message to the message type needed
	batchInfo := &pb.BatchReady{}
	err = ptypes.UnmarshalAny(authMsg.Message, batchInfo)
	if err != nil {
		return err
	}

	return s.handler.DownloadMixedBatch(stream, batchInfo, authState)
}
