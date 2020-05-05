module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200504184505-e210a219cbd9
	gitlab.com/elixxir/primitives v0.0.0-20200504203122-2adfdf770891
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/net v0.0.0-20200505041828-1ed23360d12c
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.22.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.29.1
