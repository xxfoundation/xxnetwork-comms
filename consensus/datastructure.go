package consensus

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	"sync"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000

type idFunc func(interface{}) uint64
type compFunc func(interface{}, interface{}) bool

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *RingBuff
}

func (u *Updates) AddRound(info *mixmessages.RoundInfo) error {
	// comparison should ensure that updates are not overwritten in the event of a duplicate
	comp := func(current interface{}, new interface{}) bool {
		if current == nil {
			return true
		}
		return false
	}

	// we want each updateId stored in this structure
	id := func(val interface{}) uint64 {
		return val.(*mixmessages.RoundInfo).UpdateID
	}

	if u.updates == nil {
		u.updates = New(RoundUpdatesBufLen, id)
	}

	return u.updates.UpsertById(info, id, comp)
}

// ID numbers can overwrite
type Data struct {
	rounds *RingBuff
}

func (d *Data) UpsertRound(r *mixmessages.RoundInfo) error {
	// comparison here should ensure that either the current round is nil or has a lower update id than the new round
	comp := func(current interface{}, new interface{}) bool {
		if current == nil {
			return true
		}
		if current.(*mixmessages.RoundInfo).UpdateID < new.(*mixmessages.RoundInfo).UpdateID {
			return true
		}
		return false
	}

	// We want data using the round ID as its primary
	id := func(val interface{}) uint64 {
		return val.(*mixmessages.RoundInfo).ID
	}

	if d.rounds == nil {
		d.rounds = New(RoundInfoBufLen, id)
	}

	//find the round location
	//check the new state is newer then the current
	//replace the round info object
	return d.rounds.UpsertById(r, id, comp)
}

type RingBuff struct {
	buff             []interface{}
	len, first, last int
	id               idFunc
	lock             sync.Mutex
}

// next is a helper function for ringbuff
// it handles incrementing the first & last markers
func (rb *RingBuff) next() {
	rb.last = (rb.last + 1) % rb.len
	if rb.last == rb.first {
		rb.first = (rb.first + 1) % rb.len
	}
	if rb.first == -1 {
		rb.first = 0
	}
}

// getIndex is a helper function for ringbuff
// it returns an index relative to the first/last position of the buffer
func (rb *RingBuff) getIndex(i int) int {
	var index int
	if i < 0 {
		index = (rb.last + rb.len + i) % rb.len
	} else {
		index = (rb.first + i) % rb.len
	}
	return index
}

// Initialize a new ring buffer with length n
func New(n int, id idFunc) *RingBuff {
	rb := &RingBuff{
		buff:  make([]interface{}, 0),
		len:   n,
		first: -1,
		last:  0,
		id:    id,
	}
	return rb
}

// Push a round to the buffer
func (rb *RingBuff) Push(val interface{}) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.buff[rb.last] = val
	rb.next()
}

// push a round to a relative index in the buffer
func (rb *RingBuff) UpsertById(val interface{}, id idFunc, comp compFunc) error {
	rb.lock.Lock()
	defer rb.lock.Unlock()
	newId := id(val)

	if id(rb.buff[rb.first]) > newId {
		return errors.Errorf("Did not upsert value %+v; id is older than first tracked", val)
	}

	lastId := id(rb.Get())
	if lastId+1 == newId {
		rb.Push(val)
	} else if (lastId + 1) < newId {
		for i := lastId + 1; i <= newId; i++ {
			rb.Push(nil)
		}
		rb.Push(val)
	} else if lastId+1 > newId {
		i := rb.getIndex(int(newId - (lastId + 1)))
		if comp(rb.buff[i], val) {
			rb.buff[i] = val
		} else {
			return errors.Errorf("Did not upsert value %+v; comp function returned false", val)
		}
	}
	return nil
}

func (rb *RingBuff) Get() interface{} {
	mostRecentIndex := (rb.last + rb.len - 1) % rb.len
	return rb.buff[mostRecentIndex]
}

func (rb *RingBuff) GetById(i int) (interface{}, error) {
	firstId := rb.id(rb.buff[rb.first])
	if i < int(firstId) {
		return nil, errors.Errorf("requested ID %d is lower than oldest id %d", i, firstId)
	}

	lastId := rb.id(rb.Get())
	if i > int(lastId) {
		return nil, errors.Errorf("requested id %d is higher than most recent id %d", i, lastId)
	}

	index := rb.getIndex(int(firstId) - i)
	return rb.buff[index], nil
}

func (rb *RingBuff) Len() int {
	return rb.len
}
