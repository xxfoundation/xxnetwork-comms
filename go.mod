module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200504184505-e210a219cbd9
	gitlab.com/elixxir/primitives v0.0.0-20200510221418-b069d509885c
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f
	golang.org/x/sys v0.0.0-20200509044756-6aff5f38e54f // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200511104702-f5ebc3bea380 // indirect
	google.golang.org/grpc v1.27.1
	google.golang.org/protobuf v1.22.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
