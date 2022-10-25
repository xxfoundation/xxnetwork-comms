# elixxir/comms

[![pipeline status](https://gitlab.com/elixxir/comms/badges/master/pipeline.svg)](https://gitlab.com/elixxir/comms/commits/master)
[![coverage report](https://gitlab.com/elixxir/comms/badges/master/coverage.svg)](https://gitlab.com/elixxir/comms/commits/master)

This library implements functionality for communications operations in
the xx network system.

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

First install the protobuf compiler or update by following the instructions in
[Installing Protocol Buffer Compiler](#installing-protocol-buffer-compiler)
below.

Use the following command to compile a protocol buffer.

```shell
protoc -I. -I../vendor --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto
```

* This command must be run from the directory containing the `.proto` file
  being compiled.
* The `-I` flag specifies where to find imports used by the `.proto` file and
  may need to be modified or removed to suit the .proto file being compiled.\
    * 💡 **Note:** Note: If you are importing a file from the vendor directory,
      ensure that you have the correct version by running `go mod vendor`.
* If there is more than one proto file in the directory, replace `*.proto` with
  the file’s name.
* If the `.proto` file does not use gRPC, then the `--go-grpc_out` and
  `--go-grpc_opt` can be excluded.

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

## Installing Protocol Buffer Compiler

This guide describes how to install the required dependencies to compile
`.proto` files to Go.

Before following the instructions below, be sure to remove all old versions of
`protoc`. If your previous protoc-gen-go file is not installed in your Go bin
directory, it will also need to be removed.

If you have followed this guide previously when installing `protoc` and need to
update, you can simply follow the instructions below. No uninstallation or
removal is necessary.

To compile a protocol buffer, you need the protocol buffer compiler `protoc`
along with two plugins `protoc-gen-go` and `protoc-gen-go-grpc`. Make sure you
use the correct versions as listed below.

|                      | Version | Download                                                            | Documentation                                                           |
|----------------------|--------:|---------------------------------------------------------------------|-------------------------------------------------------------------------|
| `protoc`             |  3.15.6 | https://github.com/protocolbuffers/protobuf/releases/tag/v3.15.6    | https://developers.google.com/protocol-buffers/docs/gotutorial          |
| `protoc-gen-go`      |  1.27.1 | https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.27.1 | https://pkg.go.dev/google.golang.org/protobuf@v1.27.1/cmd/protoc-gen-go |
| `protoc-gen-go-grpc` |   1.2.0 | https://github.com/grpc/grpc-go/releases/tag/v1.2.0                 | https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc        |

1. Download the correct release of `protoc` from the
   [release page](https://github.com/protocolbuffers/protobuf/releases) or use
   the link from the table above to get the download for your OS.

       wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/protoc-3.15.6-linux-x86_64.zip

2. Extract the files to a folder, such as `$HOME/.local`.

       unzip protoc-3.15.6-linux-x86_64.zip -d $HOME/.local

3. Add the selected directory to your environment’s `PATH` variable, make sure
   to include it in your `.profile` or `.bashrc` file. Also, include your go bin
   directory (`$GOPATH/bin` or `$GOBIN`) if it is not already included.

       export PATH="$PATH:$HOME/.local/bin:$GOPATH/bin"

   💡 **Note:** Make sure you update your configuration file once done with
   source `.profile`.

4. Now check that `protoc` is installed with the correct version by running the
   following command.

       protoc --version

   Which prints the current version

       libprotoc 3.15.6

5. Next, download `protoc-gen-go` and `protoc-gen-go-grpc` using the version
   found in the table above.

       go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
       go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

6. Check that `protoc-gen-go` is installed with the correct version.

       protoc-gen-go --version
       protoc-gen-go v1.27.1

7. Check that `protoc-gen-go-grpc` is installed with the correct version.

       protoc-gen-go-grpc --version
       protoc-gen-go-grpc 1.2.0