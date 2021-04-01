module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/aws/aws-lambda-go v1.8.1 // indirect
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.7-0.20210401160850-96cbf25fc59e
	gitlab.com/elixxir/primitives v0.0.3-0.20210401160752-2fe779c9fb2a
	gitlab.com/xx_network/comms v0.0.4-0.20210401160731-7b8890cdd8ad
	gitlab.com/xx_network/crypto v0.0.5-0.20210401160648-4f06cace9123
	gitlab.com/xx_network/primitives v0.0.4-0.20210331161816-ed23858bdb93
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
