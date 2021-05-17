module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.7-0.20210514185043-7398614e9e43
	gitlab.com/elixxir/primitives v0.0.3-0.20210514174939-fbba435364ed
	gitlab.com/xx_network/comms v0.0.4-0.20210507215532-38ed97bd9365
	gitlab.com/xx_network/crypto v0.0.5-0.20210504210244-9ddabbad25fd
	gitlab.com/xx_network/primitives v0.0.4-0.20210504205835-db68f11de78a
	gitlab.com/xx_network/ring v0.0.2
	golang.org/x/net v0.0.0-20201029221708-28c70e62bb1d
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
