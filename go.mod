module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/client v1.5.0
	gitlab.com/elixxir/crypto v0.0.0-20201006010428-67a8782d097e
	gitlab.com/elixxir/primitives v0.0.0-20201006010327-c2f93eb587e3
	gitlab.com/xx_network/comms v0.0.0-20200924172734-1124191b69ee
	gitlab.com/xx_network/crypto v0.0.0-20200812183430-c77a5281c686
	gitlab.com/xx_network/primitives v0.0.0-20200812183720-516a65a4a9b2
	gitlab.com/xx_network/ring v0.0.2
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.31.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
