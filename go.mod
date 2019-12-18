module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/pkg/errors v0.8.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20191121235352-86d305a9b253
	gitlab.com/elixxir/primitives v0.0.0-20191204001459-cef5ed720564
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20191218084908-4a24b4065292 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191216205247-b31c10ee225f // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
