module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200108005412-8159c60663f9
	gitlab.com/elixxir/primitives v0.0.0-20200122202948-c556f3b5ab85
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200117163144-32f20d992d24 // indirect
	google.golang.org/grpc v1.26.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
