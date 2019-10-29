module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20191028234102-337a5624a734
	gitlab.com/elixxir/primitives v0.0.0-20191028233752-882c08b8f095
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191028173616-919d9bdd9fe6 // indirect
	google.golang.org/grpc v1.24.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
