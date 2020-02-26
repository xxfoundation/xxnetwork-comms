module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/testify v1.5.1 // indirect
	gitlab.com/elixxir/crypto v0.0.0-20200206203107-b8926242da23
	gitlab.com/elixxir/primitives v0.0.0-20200218211222-4193179f359c
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d // indirect
	golang.org/x/net v0.0.0-20200225223329-5d076fcf07a8
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200225123651-fc8f55426688 // indirect
	google.golang.org/grpc v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
