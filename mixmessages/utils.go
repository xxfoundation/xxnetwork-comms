////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains utils functions for comms

package mixmessages

import (
	"github.com/golang/protobuf/proto"
	jww "github.com/spf13/jwalterweatherman"
)

// Headers for streaming
const (
	PostPhaseHeader        = "batchinfo"
	UnmixedBatchHeader     = "unmixedbatchinfo"
	MixedBatchHeader       = "mixedBatchInfo"
	FinishRealtimeHeader   = "finishRealtimeRoundInfo"
	PrecompTestBatchHeader = "precompTestBatch"
)

const NoStreamingHeaderErr = "Streaming header has no information from %s"

// ChunkSize is the size of a streaming chunk in bytes.
const ChunkSize = 1250

// ChunkHeader is the header used for by a gateway
// streaming its response for client poll. This is used for streaming
// the amount of chunks the response has been split into.
const ChunkHeader = "totalChunks"

// SplitResponseIntoChunks is a function which takes in a message and splits
// the serialized message into ChunkSize chunks. .
func SplitResponseIntoChunks(message proto.Message) ([]*StreamChunk, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	// Go will round down on integer division, the arithmetic below
	// ensures the division rounds up
	chunks := make([]*StreamChunk, 0, (len(data)+ChunkSize-1)/ChunkSize)
	for loc := 0; len(data) > loc; loc += ChunkSize {
		end := loc + ChunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, &StreamChunk{Datum: data[loc:end]})
	}

	return chunks, nil
}

// AssembleChunksIntoResponse takes a list of StreamChunk's and assembles
// the datum into the message type expected by the caller.
// This functions acts as the inverse of SplitResponseIntoChunks.
func AssembleChunksIntoResponse(chunks []*StreamChunk, response proto.Message) error {
	// Get the length of the last chunk packet
	lastChunkLen := len(chunks[len(chunks)-1].Datum)

	// Calculate the total data of all chunks
	totalData := ChunkSize*(len(chunks)-1) + lastChunkLen

	// Allocate a byte slice with that data
	data := make([]byte, 0, totalData)

	// Populate the data slice with the chunk data
	for _, chunk := range chunks {
		data = append(data, chunk.Datum...)
	}

	return proto.Unmarshal(data, response)
}

func DebugMode() {
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)
}

func TraceMode() {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)
}
