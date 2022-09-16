////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"errors"
	"gitlab.com/xx_network/primitives/id"
	"strings"
)

const baseAuthErr = "Failed to authenticate"

// AuthError returns a valid authorization error on the given id
func AuthError(id *id.ID) error {
	if id == nil {
		return errors.New(baseAuthErr + " due to nil id")
	}
	return errors.New(baseAuthErr + " id: " + id.String())
}

// IsAuthError returns true if the passed error is a valid auth error
func IsAuthError(err error) bool {
	return strings.Contains(err.Error(), baseAuthErr)
}
