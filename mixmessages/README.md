gRPC - Adding New Cryptop Message Types
----

**Note**: This guide is specifically intended for adding new `node` 
comms. It can easily generalized to comms relevant to client, gateway, etc.,
but specific file paths may not be correct.

#### Step 1: Add Message to mixmessages.proto

Create a new `message`, resembling a `struct` in golang.

```
message PrecompDecryptSlot {
  uint64 Slot = 1;
  bytes EncryptedPayloadAKeys = 2;
  bytes EncryptedPayloadBKeys = 3;
  bytes PartialPayloadACypherText = 4;
  bytes PartialPayloadBCypherText = 5;
}
```

Simply specify a type and name for each field, and set equal to its field number.

For `cryptop`-type messages, you will likely need to wrap that message in a message of its own
containing the `RoundID` and an array (aka `repeated`) of messages. This wrapper message represents
a batch for the cryptop.

```
// Message for batch of Precomp Decrypt Slots
message PrecompDecryptMessage {
  string RoundID = 1;
  repeated PrecompDecryptSlot Slots = 2;
}
```
Then, simply add an `rpc` in `service Node` specifying what the 
endpoint for your new
message will be called. You must also specify what message that endpoint will trigger with, and
what type of message to respond with (*Keep in mind, these are examples and actual names may change*)

For example, this rpc is called `PrecompDecrypt`. `PrecompDecryptMessage` messages will end up here, and
the endpoint will respond with a generic blank `Ack` message (Use `Ack` if you don't need anything back,
which is the case for cryptops).

```
rpc PrecompDecrypt (PrecompDecryptMessage) returns (Ack) {}
```

###### Optional
You may also want to create a message that should be used as an `ACK` for your new message.
For example, we want to respond to `Ping` messages with `Pong` messages:

```
rpc AskOnline (Ping) returns (Pong) {}
message Ping {}
message Pong {}
```

#### Step 2: Regenerate mixmessages.pb.go

Run the following command in the base project directory
(assuming you've set gRPC up correctly per the main README):

```
go get -u github.com/golang/protobuf/protoc-gen-go@v1.22.0
cd mixmessages
protoc -I. -I/path/to/gitlab.com mixmessages.proto --go_opt=paths=source_relative --go_out=plugins=grpc:../mixmessages/
cd ..
```

Note that `/path/to/gitlab` must have `xx_network/comms` and `elixxir/comms` checked out into it.

This enables interacting with your new messages like normal go `structs`.

#### Step 3: Implement your rpc endpoint in endpoint.go

This step is required for every `rpc` you create, otherwise there will be errors. 
Create a method on `server` for your message. It will take in a `Context` and
your message as arguments, and return the message type that was
specified in your `rpc`. For example, PrecompDecrypt looks like:

```go
func (s *server) PrecompDecrypt(ctx context.Context, msg *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompDecrypt(msg)
	return &pb.Ack{}, nil
}
```

This method will be called every time the endpoint receives your message. For cryptops,
this is where you must pass your message to `server` by calling `serverHandler.YOURCRYPTOP(msg)`
like above. We will create this interface method in Step 5. 

#### Step 4: Add SendMessage function for your rpc in node package

Create a new file in `node` package for your `rpc`. The purpose of this
is to be able to send a message from the `server` repository without any dependencies on gRPC.
You may copy one of the other files in this package and modify the message input and return types.
Additionally, make sure you call your endpoint from Step 3 in the method body as follows:

`result, err := c.PrecompDecrypt(ctx, input)`

Add any additional logic that may be required when sending your message here.
This is the last step for normal messages. For cryptop messages, continue to Step 5.

#### Step 5: Add interface method for cryptop in serverHandler.go

Add a method to the interface for your new cryptop message. For example,

```go
type ServerHandler interface {
	// Server Interface for the PrecompDecrypt Messages
	PrecompDecrypt(*mixmessages.PrecompDecryptMessage)
}
```

We will be implementing this method in the `server` repository in the `server/node` package.
It is recommended that you stub this method out now in order to prevent interface implementation
errors once your new message is merged.

#### Step 6: Testing

Find the `mockserver_test.go` file in the same `node` package, and add a
blank method in order to implement the interface method you added in Step 5. Do
the same for `mockserver_test.go` located in the `client` package.

 Then, you may write a test for your `Send` function you added in Step 4 (which
 also tests the Step 3 endpoint, which is why we need the `TestInterface`). For
 example:

```go
// Smoke test SendPrecompDecrypt
func TestSendPrecompDecrypt(t *testing.T) {
	addr := "localhost:5555"
	go mixserver.StartServer(addr, TestInterface{})
	_, err := SendPrecompDecrypt(addr, &pb.PrecompDecryptMessage{})
	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
}
```
