#!/bin/bash

#///////////////////////////////////////////////////////////////////////////////
#/ Copyright © 2020 xx network SEZC                                           //
#/                                                                            //
#/ Use of this source code is governed by a license that can be found in the  //
#/ LICENSE file                                                               //
#///////////////////////////////////////////////////////////////////////////////

protoc -I. -I../vendor/ mixmessages.proto  --go_opt=paths=source_relative --go_out=plugins=grpc:../mixmessages/