////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains notificationBot -> all servers functionality

package notificationBot

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
)

// NotificationBot -> Permissioning
// Fixme: figure out what to do with notification bot and unified polling
func (nb *Comms) RequestNdf(host *connect.Host, message *pb.NDFHash) (*pb.NDF, error) {

	// Call the ProtoComms RequestNdf call
	return nb.RequestNdf(host, message)
}
