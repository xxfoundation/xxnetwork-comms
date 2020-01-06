module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20191121235352-86d305a9b253
	gitlab.com/elixxir/primitives v0.0.0-20200106183011-a68f1e6f188e
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20200106162015-b016eb3dc98e // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191230161307-f3c370f40bfb // indirect
	google.golang.org/grpc v1.26.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
