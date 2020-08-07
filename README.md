# xx_network/comms

[![pipeline status](https://gitlab.com/xx_network/comms/badges/master/pipeline.svg)](https://gitlab.com/xx_network/comms/commits/master)
[![coverage report](https://gitlab.com/xx_network/comms/badges/master/coverage.svg)](https://gitlab.com/xx_network/comms/commits/master)

#### Generating messages.pb.go


Run the following command in the base project directory
(assuming you've set gRPC up correctly):

```
protoc -I messages/ messages/messages.proto -I /path/to/gitlab.com/ --go_out=plugins=grpc:messages
```

* NOTE: For `-I /path/to/gitlab.com/` in the above command, you will want to exclude 
`gitlab.com` from the path. If you wish, you may simply put 
`import google/protobuf/any.proto` in `messages.proto` and exclude `I /path/to/gitlab.com/` entirely.
This structure allows for compatibility with our partner organization, Elixxir and their projects.

```
protoc -I messages/ messages/messages.proto --go_out=plugins=grpc:messages
protoc -I gossip/ gossip/messages.proto --go_out=plugins=grpc:gossip
```
