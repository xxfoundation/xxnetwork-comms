module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200504184505-e210a219cbd9
	gitlab.com/elixxir/primitives v0.0.0-20200506184657-93da24058321
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79 // indirect
	golang.org/x/net v0.0.0-20200506145744-7e3656a0809f
	golang.org/x/sys v0.0.0-20200508214444-3aab700007d7 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200507105951-43844f6eee31 // indirect
	google.golang.org/grpc v1.27.1
	google.golang.org/protobuf v1.22.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
