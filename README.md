# xx_network/comms

[![pipeline status](https://gitlab.com/xx_network/comms/badges/master/pipeline.svg)](https://gitlab.com/xx_network/comms/commits/master)
[![coverage report](https://gitlab.com/xx_network/comms/badges/master/coverage.svg)](https://gitlab.com/xx_network/comms/commits/master)

```
protoc -I messages/ messages/messages.proto -I /path/to/gitlab.com/ --go_out=plugins=grpc:messages
```