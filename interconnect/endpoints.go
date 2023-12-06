////////////////////////////////////////////////////////////////////////////////
// Copyright © 2024 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Endpoints for the interconnect service

package interconnect

import (
	"context"
	"gitlab.com/xx_network/comms/messages"
)

func (c *CMixServer) GetNDF(ctx context.Context, ping *messages.Ping) (*NDF, error) {
	return c.handler.GetNDF()
}
