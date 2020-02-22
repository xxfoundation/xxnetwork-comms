module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200221215124-b2ae9abd2bd0
	gitlab.com/elixxir/primitives v0.0.0-20200218211222-4193179f359c
	golang.org/x/net v0.0.0-20200219183655-46282727080f
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200218151345-dad8c97a84f5 // indirect
	google.golang.org/grpc v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
