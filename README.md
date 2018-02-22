# privategrity/comms

[![pipeline status](https://gitlab.com/privategrity/comms/badges/master/pipeline.svg)](https://gitlab.com/privategrity/comms/commits/master)
[![coverage report](https://gitlab.com/privategrity/comms/badges/master/coverage.svg)](https://gitlab.com/privategrity/comms/commits/master)

This library implements functionality for communications operations in
the cMix system.

## Regenerate Protobuf File

To regenerate the `mixmessage.pb.go` file, first install gRPC:
`
go get -u google.golang.org/grpc
`

Then install protocol buffers v3 as follows:
- The simplest way to do this is to download pre-compiled binaries 
  **for your platform** (`protoc-[version]-[platform].zip`) from here: https://github.com/google/protobuf/releases
	- Unzip and move the protoc binary in the `bin` directory to some logical location
	  (e.g. `/usr/local/bin` on Mac OS X)
	- Update the environment variable `PATH` to include the path to the protoc binary file you moved as follows:
		- Find your `.bash_profile` file and add the following line (if it isn't already there) with the path set to
		  wherever you placed the protoc binary: `export PATH=/usr/local/bin:$PATH`

Next, install the protoc plugin for go with the following command:

`
go get -u github.com/golang/protobuf/protoc-gen-go
`

 This will add the plugin to your `go` directory in a `bin` folder. You must add this to your `PATH` variable
 in your `.bash_profile`, so again find your `.bash_profile file` and make sure your `GOPATH` is defined
 (on Mac OS the default is: `export GOPATH=$HOME/go`).
 
 If `GOPATH` isn't defined, then add a line like the above defining it (replacing `$HOME/go` with the
 path of your go directory). Finally add the following line to your `.bash_profile`:
 
`
export PATH=$PATH:$GOPATH/bin
`



Everything should be installed now. Check by running `protoc-gen-go` in a terminal. It should hang with no errors
if everything is installed correctly. Now navigate to the `comms` project directory and run the following command 
in the terminal in order to regenerate the `mixmessage.pb.go` file:

`
protoc -I mixmessages/ mixmessages/mixmessages.proto --go_out=plugins=grpc:mixmessages
`
