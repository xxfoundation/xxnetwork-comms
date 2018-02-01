GRPC - Adding New Cryptop Message Types
----

#### Step 1: Add Message to mixmessages.proto

Create a new `message`, resembling a `struct` in golang.

```golang
message PrecompDecryptSlot {
  uint64 Slot = 1;
  bytes EncryptedMessageKeys = 2;
  bytes EncryptedRecipientIDKeys = 3;
  bytes PartialMessageCypherText = 4;
  bytes PartialRecipientIDCypherText = 5;
}
```

Simply specifiy a type and name for each field, and set equal to its field number.

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
Then, simply add an `rpc` in `service MixMessageService` specifying what the endpoint for your new
message will be called. You must also specify what message that endpoint will trigger with, and
what type of message to respond with.

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
(assuming you've set GRPC up correctly per the main README):

`protoc -I mixmessages/ mixmessages/mixmessages.proto --go_out=plugins=grpc:mixmessages`

This enables interacting with your new messages like normal go `structs`.

#### Step 3: Implement your rpc endpoint in mixserver.go
