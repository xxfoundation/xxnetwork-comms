package connect

import (
	"errors"
	"strings"
)

const baseAuthErr = "Failed to authenticate"

// AuthError returns a valid authorization error on the given id
func AuthError(id string) error {
	return errors.New(baseAuthErr + " id: " + id)
}

// IsAuthError returns true if the passed error is a valid auth error
func IsAuthError(err error) bool {
	return strings.Contains(err.Error(), baseAuthErr)
}