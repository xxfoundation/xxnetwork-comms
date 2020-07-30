package network

import (
	"errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
)

func (i *Instance) GetHistoricalRound(id id.Round) (*pb.RoundInfo, error) {
	if i.ers != nil {
		return i.ers.Retrieve(id)
	}
	return nil, errors.New("no ExternalRoundStorage object was defined on instance creation")
}

func (i *Instance) GetHistoricalRounds(rounds []id.Round) ([]*pb.RoundInfo, error) {
	if i.ers != nil {
		return i.ers.RetrieveMany(rounds)
	}
	return nil, errors.New("no ExternalRoundStorage object was defined on instance creation")
}

func (i *Instance) GetHistoricalRoundRange(first, last id.Round) ([]*pb.RoundInfo, error) {
	if i.ers != nil {
		return i.ers.RetrieveRange(first, last)
	}
	return nil, errors.New("no ExternalRoundStorage object was defined on instance creation")
}
