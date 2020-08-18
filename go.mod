module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/testify v1.6.1 // indirect
	gitlab.com/elixxir/crypto v0.0.0-20200707005343-97f868cbd930
	gitlab.com/elixxir/primitives v0.0.0-20200706165052-9fe7a4fb99a3
	gitlab.com/xx_network/comms v0.0.0-20200730220144-eea32e8b696d
	golang.org/x/crypto v0.0.0-20200707235045-ab33eee955e0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200709005830-7a2ca40e9dc3 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
