package consensus

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"sync"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000


type Updates struct{
	updates [RoundUpdatesBufLen]*mixmessages.RoundInfo
}


type Data struct{
	rounds [RoundInfoBufLen]*mixmessages.RoundInfo
	earliestRound uint64
	letestRound uint64
}

func (d *Data)UpsertRound(r *mixmessages.RoundInfo)(error){
	if d.earliestRound> r.ID {
		return errors.New("update to untracked round")
	}

	//find the round location
	//check the new state is newer then the current
	//replace the round info object
}







