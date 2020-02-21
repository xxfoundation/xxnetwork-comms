module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/testify v1.5.1 // indirect
	gitlab.com/elixxir/crypto v0.0.0-20200221004027-c949af0228ec
	gitlab.com/elixxir/primitives v0.0.0-20200218211222-4193179f359c
	golang.org/x/crypto v0.0.0-20200221170553-0f24fbd83dfb // indirect
	golang.org/x/net v0.0.0-20200219183655-46282727080f
	golang.org/x/sys v0.0.0-20200219091948-cb0a6d8edb6c // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200218151345-dad8c97a84f5 // indirect
	google.golang.org/grpc v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
