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
	"testing"
)

func TestAuthError(t *testing.T) {
	expectedAuthErrorStr := "Failed to authenticate id: c29pc29pc29pAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	result := AuthError(id.NewIdFromString("soisoisoi", id.Generic, t))
	if result == nil {
		t.Error("AuthError did not return an error object")
	}
	if result.Error() != expectedAuthErrorStr {
		t.Errorf("returned error not as expected: Expected: %s, received: %s",
			expectedAuthErrorStr, result.Error())
	}

	result = AuthError(nil)
	if result == nil || !strings.Contains(result.Error(), "due to nil id") {
		t.Errorf("returned error not as expected: Expected: %s, received: %s",
			"due to nil id", result.Error())
	}

}

func TestIsAuthError(t *testing.T) {
	isAuthError := errors.New("Failed to authenticate id: soisoisoi")

	if !IsAuthError(isAuthError) {
		t.Errorf("IsAuthError returned that authError is not an authError")
	}

	notAuthError := errors.New("dont feed the plants")

	if IsAuthError(notAuthError) {
		t.Errorf("IsAuthError returned that a non authError is an authError")
	}
}
