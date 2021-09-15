# xx_network/comms

[![pipeline status](https://gitlab.com/xx_network/comms/badges/master/pipeline.svg)](https://gitlab.com/xx_network/comms/commits/master)
[![coverage report](https://gitlab.com/xx_network/comms/badges/master/coverage.svg)](https://gitlab.com/xx_network/comms/commits/master)

#### Generating messages.pb.go

First, ensure protoc is installed and get 1.22.0 (as of this writing) of
protoc-gen-go:

```
go get -u github.com/golang/protobuf/protoc-gen-go
```


Run the following command in the base project directory
(assuming you've set gRPC up correctly):

```
cd messages
protoc -I. messages.proto --go_opt=paths=source_relative --go_out=plugins=grpc:../messages/
cd ..
```

Note that `/path/to/gitlab.com` is required to generate correct import lines.
This lib needs to be checked out inside of the
`/path/to/gitlab.com/xx_network/comms` folder.

```
cd interconnect
protoc -I. interconnect.proto -I/path/to/gitlab.com --go_opt=paths=source_relative --go_out=plugins=grpc:../interconnect/
cd ../gossip
protoc -I. gossip.proto --go_opt=paths=source_relative --go_out=plugins=grpc:../gossip/
cd ..
```
