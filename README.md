# xx_network/comms

[![pipeline status](https://gitlab.com/xx_network/comms/badges/master/pipeline.svg)](https://gitlab.com/xx_network/comms/commits/master)
[![coverage report](https://gitlab.com/xx_network/comms/badges/master/coverage.svg)](https://gitlab.com/xx_network/comms/commits/master)

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
| `protoc`             |  3.21.9 | https://github.com/protocolbuffers/protobuf/releases/tag/v3.21.9    | https://developers.google.com/protocol-buffers/docs/gotutorial          |
| `protoc-gen-go`      |  1.28.1 | https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.28.1 | https://pkg.go.dev/google.golang.org/protobuf@v1.28.1/cmd/protoc-gen-go |
| `protoc-gen-go-grpc` |   1.2.0 | https://github.com/grpc/grpc-go/releases/tag/v1.2.0                 | https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc        |

1. Download the correct release of `protoc` from the
   [release page](https://github.com/protocolbuffers/protobuf/releases) or use
   the link from the table above to get the download for your OS.

       wget https://github.com/protocolbuffers/protobuf/releases/download/v3.21.9/protoc-3.21.9-linux-x86_64.zip

2. Extract the files to a folder, such as `$HOME/.local`.

       unzip protoc-3.21.9-linux-x86_64.zip -d $HOME/.local

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

       libprotoc 3.21.9

5. Next, download `protoc-gen-go` and `protoc-gen-go-grpc` using the version
   found in the table above.

       go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
       go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

6. Check that `protoc-gen-go` is installed with the correct version.

       protoc-gen-go --version
       protoc-gen-go v1.28.1

7. Check that `protoc-gen-go-grpc` is installed with the correct version.

       protoc-gen-go-grpc --version
       protoc-gen-go-grpc 1.2.0