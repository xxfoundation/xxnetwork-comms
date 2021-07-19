package node

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc/metadata"
)

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
