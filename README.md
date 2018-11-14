# elixxir/comms

[![pipeline status](https://gitlab.com/elixxir/comms/badges/master/pipeline.svg)](https://gitlab.com/elixxir/comms/commits/master)
[![coverage report](https://gitlab.com/elixxir/comms/badges/master/coverage.svg)](https://gitlab.com/elixxir/comms/commits/master)

This library implements functionality for communications operations in
the cMix system.

## How to run tests

First, make sure dependencies are installed into the vendor folder by running
`glide up`. Then, in the project directory, run `go test ./...`.

If what you're working on requires you to change other repos, you can remove
the other repo from the vendor folder and Go's build tools will look for those
packages in your Go path instead. Knowing which dependencies to remove can be
really helpful if you're changing a lot of repos at once.

If glide isn't working and you don't know why, try removing glide.lock and
~/.glide to brutally cleanse the cache.

## Regenerate Protobuf File

To regenerate the `mixmessage.pb.go` file, first install gRPC:

`go get -u google.golang.org/grpc`

Then install protocol buffers v3 as follows:
- The simplest way to do this is to download pre-compiled binaries
  **for your platform** (`protoc-[version]-[platform].zip`) from here:
    https://github.com/google/protobuf/releases
- Unzip and move the protoc binary in the `bin` directory to some logical
  location (e.g. `/usr/local/bin` on Mac OS X)
- Update the environment variable `PATH` to include the path to the protoc
  binary file you moved as follows:
- Find your `.bash_profile` file and add the following line (if it isn't
  already there) with the path set to wherever you placed the protoc binary:
  `export PATH=/usr/local/bin:$PATH`

Next, install the protoc plugin for go with the following command:

`go get -u github.com/golang/protobuf/protoc-gen-go`

This will add the plugin to your `go` directory in a `bin` folder. You
must add this to your `PATH` variable in your `.bash_profile`, so
again find your `.bash_profile file` and make sure your `GOPATH` is
defined (on Mac OS the default is: `export GOPATH=$HOME/go`).

If `GOPATH` isn't defined, then add a line like the above defining it
(replacing `$HOME/go` with the path of your go directory). Finally add
the following line to your `.bash_profile`:

`export PATH=$PATH:$GOPATH/bin`

Everything should be installed now. Check by running `protoc-gen-go`
in a terminal. It should hang with no errors if everything is
installed correctly. Now navigate to the `comms` project directory and
run the following command in the terminal in order to regenerate the
`mixmessage.pb.go` file:

```
protoc -I mixmessages/ mixmessages/mixmessages.proto --go_out=plugins=grpc:mixmessages
```

## Repository Organization

This repository is organized into 4 key folders:
1. `mixmessages` - contains the gRPC proto spec file and the file generated from
   that, along with any future shared helper functionality.
2. `client` - client functions for proper cmix clients.
3. `node` - gRPC endpoints hosted by cMix servers.
4. `gateway` - gateway endpoints and functions

Note that `gateway` and `node` are organized similarly. The `endpoints.go` file contains
gRPC endpoint implementations, and the `handler.go` file contains the handler interface
used by the endpoint implementations as well as the implementation struct.

## Adding Network Endpoints

To add an endpoint, you need to make changes to `handler.go` and `endpoint.go`. You will
also need to add a client function to call your new endpoint in the repo where the client
is implemented.

`handler.go`:
1. Add your function to the interface.
2. Add it to the `implementationFunctions` struct.
3. Add the "Unimplemented warning" version of the function to `NewImplementation()`
4. Add the wrapper call to `s.Functions.FUNCNAME(...)` at the bottom to implement the 
   interface for Implementation struct.

`endpoint.go`:
1. Add the gRPC implementation which calls the function through the handler.

`Client Function`:
1. Implement the client call either via a new module or an existing module. It should
   go in the location where the module is the client (i.e., if the node calls to the
   gateway, it goes in node)

