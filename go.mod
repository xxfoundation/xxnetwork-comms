module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200514211323-de119c132b78
	gitlab.com/elixxir/primitives v0.0.0-20200514181428-14736275e533
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120
	golang.org/x/sys v0.0.0-20200513112337-417ce2331b5c // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200514193133-8feb7f20f2a2 // indirect
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
