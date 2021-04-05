module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/aws/aws-lambda-go v1.8.1 // indirect
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.7-0.20210405224356-e2748985102a
	gitlab.com/elixxir/primitives v0.0.3-0.20210405224302-03092a268566
	gitlab.com/xx_network/comms v0.0.4-0.20210405224241-5447394f79d7
	gitlab.com/xx_network/crypto v0.0.5-0.20210405224157-2b1f387b42c1
	gitlab.com/xx_network/primitives v0.0.4-0.20210402222416-37c1c4d3fac4
	gitlab.com/xx_network/ring v0.0.2
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/net v0.0.0-20201029221708-28c70e62bb1d
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
