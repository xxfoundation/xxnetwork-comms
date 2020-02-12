module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200206203107-b8926242da23
	gitlab.com/elixxir/primitives v0.0.0-20200210205543-5c55c1f6949f
	golang.org/x/crypto v0.0.0-20200210222208-86ce3cb69678 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200211111953-2dc5924e3898 // indirect
	google.golang.org/grpc v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
