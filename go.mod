module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.3.4
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200229000841-b1ee7117a1d0
	gitlab.com/elixxir/primitives v0.0.0-20200305204536-da5cbdcd2d96
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200305110556-506484158171 // indirect
	google.golang.org/grpc v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
