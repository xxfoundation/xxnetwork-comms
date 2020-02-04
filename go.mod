module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200108005412-8159c60663f9
	gitlab.com/elixxir/primitives v0.0.0-20200131183153-e93c6b75019f
	golang.org/x/crypto v0.0.0-20200204104054-c9f3fb736b72 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200204135345-fa8e72b47b90 // indirect
	google.golang.org/grpc v1.27.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
