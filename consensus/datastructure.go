package consensus

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/mixmessages"
	"sync"
)

const RoundInfoBufLen = 1000
const RoundUpdatesBufLen = 10000

// Standard ring buffer, but objects come with numbering
type Updates struct {
	updates *RingBuff
}

func (u *Updates) AddRound(info *mixmessages.RoundInfo) {
	if u.updates == nil {
		u.updates = New(RoundUpdatesBufLen)
	}

	nextId := u.updates.GetMostRecent() + 1
	if nextId == info.ID {
		u.updates.Push(info)
	} else if nextId < info.ID {
		for i := nextId; i < info.ID; i++ {
			u.updates.Push(nil)
		}
		u.updates.Push(info)
	} else if nextId > info.ID {
		i := int(info.ID - nextId)
		cur := u.updates.Get(i)
		if cur != nil {
			u.updates.PushToIndex(info, i)
		} else {
			// what happens when rounds collide
		}
	}
}

// ID numbers can overwrite
type Data struct {
	rounds        [RoundInfoBufLen]*mixmessages.RoundInfo
	earliestRound uint64
	letestRound   uint64
}

func (d *Data) UpsertRound(r *mixmessages.RoundInfo) error {
	if d.earliestRound > r.ID {
		return errors.New("update to untracked round")
	}

	//find the round location
	//check the new state is newer then the current
	//replace the round info object

}

type RingBuff struct {
	buff             []interface{}
	len, first, last int
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
func New(n int) *RingBuff {
	rb := &RingBuff{
		buff:  make([]interface{}, 0),
		len:   n,
		first: -1,
		last:  0,
	}
	return rb
}

// Push a round to the buffer
func (rb *RingBuff) Push(r *mixmessages.RoundInfo) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	rb.buff[rb.last] = r
	rb.next()
}

// push a round to a relative index in the buffer
func (rb *RingBuff) PushToIndex(r *mixmessages.RoundInfo, i int) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	index := rb.getIndex(i)
	rb.buff[index] = r
}

func (rb *RingBuff) GetMostRecent() interface{} {
	mostRecentIndex := (rb.last + rb.len - 1) % rb.len
	return rb.buff[mostRecentIndex]
}

func (rb *RingBuff) Get(i int) interface{} {
	index := rb.getIndex(i)
	return rb.buff[index]
}

func (rb *RingBuff) Len() int {
	return rb.len
}
