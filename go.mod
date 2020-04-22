module gitlab.com/elixxir/comms

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/crypto v0.0.0-20200410231849-90e859940f5d
	gitlab.com/elixxir/primitives v0.0.0-20200410231944-a57d71d577c9
	golang.org/x/crypto v0.0.0-20200420201142-3c4aac89819a // indirect
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200420144010-e5e8543f8aeb // indirect
	google.golang.org/grpc v1.29.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
