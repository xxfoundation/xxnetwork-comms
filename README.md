# xx_network/comms

[![pipeline status](https://gitlab.com/xx_network/comms/badges/master/pipeline.svg)](https://gitlab.com/xx_network/comms/commits/master)
[![coverage report](https://gitlab.com/xx_network/comms/badges/master/coverage.svg)](https://gitlab.com/xx_network/comms/commits/master)

#### Generating messages.pb.go


Run the following command in the base project directory
(assuming you've set gRPC up correctly):

```
protoc -I messages/ messages/messages.proto --go_out=plugins=grpc:messages
```

```
protoc -I interconnect/ interconnect/interconnect.proto -I /path/to/gitlab.com/ --go_out=plugins=grpc:interconnect
protoc -I messages/ messages/messages.proto --go_out=plugins=grpc:messages
protoc -I gossip/ gossip/messages.proto --go_out=plugins=grpc:gossip
```
