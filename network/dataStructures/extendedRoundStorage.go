package dataStructures

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
)

type ExternalRoundStorage interface {
	// Store: stores the round info inside the underlying storage medium, which generally is a database. Store will
	// add the round info to the database if it doesn't exist and will only overwrite the data if it does exist in the
	// event that the update ID of the passed in data is greater than the update ID of the existing round info.
	Store(*pb.RoundInfo) error
	// Retrieve will return the round info for the given round ID and will return nil but not an error if it does not
	// exist.
	Retrieve(id id.Round) (*pb.RoundInfo, error)
	// RetrieveMany will return all rounds passed in in the ID list, if the round doesn't its reciprocal entry in the
	// returned slice will be blank.
	RetrieveMany(rounds []id.Round) ([]*pb.RoundInfo, error)
	// RetrieveRange will return all rounds in the range, if the round doesn't exist the reciprocal entry in the
	// returned slice will be blank.
	RetrieveRange(first, last id.Round) ([]*pb.RoundInfo, error)
}
