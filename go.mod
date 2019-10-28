module gitlab.com/elixxir/comms

go 1.12

require (
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20191024163612-1d67595237e4
	gitlab.com/elixxir/primitives v0.0.0-20191024163559-539ed465e76f
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094
	golang.org/x/sys v0.0.0-20191024172528-b4ff53e7a1cb // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191009194640-548a555dbc03 // indirect
	google.golang.org/grpc v1.24.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
