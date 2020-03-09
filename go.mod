module gitlab.com/elixxir/comms

go 1.13

require (
	cloud.google.com/go/pubsub v1.3.0 // indirect
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20191203043605-d42048ed14fd // indirect
	github.com/cncf/udpa/go v0.0.0-20200124205748-db4b343e48c1 // indirect
	github.com/envoyproxy/go-control-plane v0.9.4 // indirect
	github.com/golang/protobuf v1.3.4
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/stretchr/objx v0.2.0 // indirect
	gitlab.com/elixxir/crypto v0.0.0-20200229000841-b1ee7117a1d0
	gitlab.com/elixxir/primitives v0.0.0-20200306214728-35300c4b5152
	golang.org/x/exp v0.0.0-20200228211341-fcea875c7e85 // indirect
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	golang.org/x/mobile v0.0.0-20200222142934-3c8601c510d0 // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/tools v0.0.0-20200309162502-c94e1fe1450c // indirect
	google.golang.org/genproto v0.0.0-20200309141739-5b75447e413d // indirect
	google.golang.org/grpc v1.27.1
	rsc.io/sampler v1.99.99 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
