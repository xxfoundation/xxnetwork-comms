#!/bin/bash

#///////////////////////////////////////////////////////////////////////////////
#/ Copyright © 2020 xx network SEZC                                           //
#/                                                                            //
#/ Use of this source code is governed by a license that can be found in the  //
#/ LICENSE file                                                               //
#///////////////////////////////////////////////////////////////////////////////

protoc -I. -I../vendor --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto