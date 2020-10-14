module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/testify v1.6.1 // indirect
	gitlab.com/elixxir/crypto v0.0.0-20201002151041-c4ab8f8033dc
	gitlab.com/elixxir/primitives v0.0.0-20200930214918-50b3c2030f26
	gitlab.com/xx_network/comms v0.0.0-20200925191822-08c0799a24a6
	gitlab.com/xx_network/crypto v0.0.0-20200812183430-c77a5281c686
	gitlab.com/xx_network/primitives v0.0.0-20200812183720-516a65a4a9b2
	gitlab.com/xx_network/ring v0.0.2
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb
	golang.org/x/sys v0.0.0-20201013132646-2da7054afaeb // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20201013134114-7f9ee70cb474 // indirect
	google.golang.org/grpc v1.33.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
