///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Endpoints for the interconnect service

package interconnect

import (
	"context"
	"gitlab.com/xx_network/comms/messages"
)

func (c *Server) GetNDF(ctx context.Context, ping *messages.Ping) (*NDF, error) {
	return c.handler.GetNDF()
}
