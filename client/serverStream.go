///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	"context"
	"github.com/ktr0731/grpc-web-go-client/grpcweb"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc/metadata"
	"io"
	"strings"
)

var _ pb.Gateway_PollClient = (*serverStream)(nil)

// serverStream wraps a grpcweb.ServerStream to make it adhere to the
// pb.Gateway_PollClient interface.
type serverStream struct {
	ctx context.Context
	grpcweb.ServerStream
}

func newServerStream(
	ctx context.Context, ws grpcweb.ServerStream) pb.Gateway_PollClient {
	return &serverStream{ctx, ws}
}

func (s *serverStream) Recv() (*pb.StreamChunk, error) {
	m := new(pb.StreamChunk)
	if err := s.ServerStream.Receive(s.ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *serverStream) Header() (metadata.MD, error) {
	return s.ServerStream.Header()
}

func (s *serverStream) Trailer() metadata.MD {
	return s.ServerStream.Trailer()
}

func (s *serverStream) CloseSend() error {
	return nil
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

func (s *serverStream) SendMsg(m interface{}) error {
	return s.ServerStream.Send(s.ctx, m)
}

func (s *serverStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.Receive(s.ctx, m)

	// If the response has been closed, return EOF
	const responseClosedErr = "http: read on closed response body"
	if err != nil && strings.Contains(err.Error(), responseClosedErr) {
		return io.EOF
	}

	return err
}
