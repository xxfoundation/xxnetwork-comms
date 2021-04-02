package mixmessages

import (
	"github.com/golang/protobuf/ptypes"
	"math/rand"
	"testing"
	"time"
)

func Test_TestMessages(t *testing.T) {
	lengthOfMessageSlice := 500

	// Construct messages
	oneKbMessage := populateMessages(lengthOfMessageSlice, t, makeOneKilobyte)
	twoKbMessage := populateMessages(lengthOfMessageSlice, t, makeTwoKilobyte)
	nilMessage := populateMessages(lengthOfMessageSlice, t, makeNilMessage)

	// Construct any message type for each message above
	anyOneKbMsg, err := ptypes.MarshalAny(oneKbMessage)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}

	anyTwoKbMsg, err := ptypes.MarshalAny(twoKbMessage)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}


	anyNilMsg, err := ptypes.MarshalAny(nilMessage)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}


	// Print
	t.Logf("message length: %d", len(anyOneKbMsg.Value))
	t.Logf("message length: %d", len(anyTwoKbMsg.Value))
	t.Logf("message length: %d", len(anyNilMsg.Value))


}

func populateMessages(len int, t *testing.T,
	populateMessage func() ([]byte, int)) *TestMessage {
	testMessage := &TestMessage{
		Data:      make([][]byte, len),
	}
	counter := 0
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)
	var dataLen int
	for i := 0; i < len; i++ {
		if r.Int() % 2 == 0 {
			counter++
			testMessage.Data[i], dataLen = populateMessage()
		}
	}
	
	t.Logf("Populated %d of %d fields with %dKB of data", counter, len, dataLen/1000)

	return testMessage
}

func makeNilMessage() ([]byte, int) {
	byteLen := 0
	data := make([]byte, byteLen)
	for i := 0; i < byteLen; i ++{
		data[i] = byte(i)
	}

	return data, 0

}


func makeOneKilobyte() ([]byte, int) {
	byteLen := 1000
	data := make([]byte, byteLen)
	for i := 0; i < byteLen; i ++{
		data[i] = byte(i)
	}

	return data, byteLen
}

func makeTwoKilobyte() ([]byte, int) {
	byteLen := 2000
	data := make([]byte, byteLen)
	for i := 0; i < byteLen; i ++{
		data[i] = byte(i)
	}

	return data, byteLen
}