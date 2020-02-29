module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200206203107-b8926242da23
	gitlab.com/elixxir/primitives v0.0.0-20200227004200-0eb9f2db7d5d
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200228133532-8c2c7df3a383 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
