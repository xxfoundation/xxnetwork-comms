module gitlab.com/elixxir/comms

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20191025205058-be1aca8c7a70
	gitlab.com/elixxir/primitives v0.0.0-20191025204417-8f0ead5b2d6b
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03 // indirect
	google.golang.org/grpc v1.24.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
