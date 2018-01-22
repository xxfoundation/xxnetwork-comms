# Comms Lib

This library implements functionality for communications operations in
the cMix system.

## Regenerate Protobuf File

To regenerate the `mixmessage.pb.go` file, run the following command:


```
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```
