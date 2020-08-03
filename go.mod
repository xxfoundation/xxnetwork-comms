module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/testify v1.6.1 // indirect
	gitlab.com/elixxir/crypto v0.0.0-20200721213839-b026955c55c0
	gitlab.com/elixxir/primitives v0.0.0-20200731184040-494269b53b4d
	gitlab.com/xx_network/collections/ring v0.0.0-00010101000000-000000000000
	gitlab.com/xx_network/comms v0.0.0-20200731231107-9e020daf0013
	golang.org/x/crypto v0.0.0-20200707235045-ab33eee955e0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200709005830-7a2ca40e9dc3 // indirect
	google.golang.org/grpc v1.30.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace (
	gitlab.com/xx_network/collections/ring => gitlab.com/xx_network/collections/ring.git v0.0.1
	google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
)
