module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200513205206-8be446a9ccbe
	gitlab.com/elixxir/primitives v0.0.0-20200514163806-e7b4d77801d8
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120
	golang.org/x/sys v0.0.0-20200513112337-417ce2331b5c // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884 // indirect
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
