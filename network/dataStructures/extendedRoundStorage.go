////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/primitives/id"
)

// The ExtendedRoundStorage (ERS) interface allows storing rounds inside of an external database for clients to pull
// from, because the ring buffer only contains a limited number of them while clients might need to go further back
// into history.
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
