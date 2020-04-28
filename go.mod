module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200410231849-90e859940f5d
	gitlab.com/elixxir/primitives v0.0.0-20200428170743-1263bbc7df36
	golang.org/x/crypto v0.0.0-20200427165652-729f1e841bcc // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	golang.org/x/sys v0.0.0-20200427175716-29b57079015a // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200428115010-c45acf45369a // indirect
	google.golang.org/grpc v1.29.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
