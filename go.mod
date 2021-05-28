module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.2
	github.com/nyaruka/phonenumbers v1.0.60 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.7-0.20210526002540-1fb51df5b4b2
	gitlab.com/elixxir/primitives v0.0.3-0.20210526002350-b9c947fec050
	gitlab.com/xx_network/comms v0.0.4-0.20210528015510-2b219e8fdfc7
	gitlab.com/xx_network/crypto v0.0.5-0.20210526002149-9c08ccb202be
	gitlab.com/xx_network/primitives v0.0.4-0.20210525232109-3f99a04adcfd
	gitlab.com/xx_network/ring v0.0.3-0.20210527191221-ce3f170aabd5
	golang.org/x/net v0.0.0-20201029221708-28c70e62bb1d
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
