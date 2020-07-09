module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200707005343-97f868cbd930
	gitlab.com/elixxir/primitives v0.0.0-20200706165052-9fe7a4fb99a3
	gitlab.com/xx_network/comms v0.0.0-20200709165104-1fcde4b1729d
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.30.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
