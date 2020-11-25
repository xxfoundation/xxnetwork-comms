module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/zeebo/blake3 v0.1.0 // indirect
	github.com/zeebo/pcg v1.0.0 // indirect
	gitlab.com/elixxir/crypto v0.0.5-0.20201125005724-bcc603df02d3
	gitlab.com/elixxir/primitives v0.0.3-0.20201116174806-97f190989704
	gitlab.com/xx_network/comms v0.0.4-0.20201119231004-a67d08045535
	gitlab.com/xx_network/crypto v0.0.5-0.20201124194022-366c10b1bce0
	gitlab.com/xx_network/primitives v0.0.3-0.20201116234927-44e42fc91e7c
	gitlab.com/xx_network/ring v0.0.2
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/net v0.0.0-20201029221708-28c70e62bb1d
	golang.org/x/sys v0.0.0-20201029080932-201ba4db2418 // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20201030142918-24207fddd1c3 // indirect
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
