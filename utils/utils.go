package utils

import (
	"github.com/mitchellh/go-homedir"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
)

// Given a path, replace a "~" character
// with the home directory to return a full file path
func GetFullPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			jww.FATAL.Panicf("Unable to locate home directory: %v", err)
		}
		// Append the home directory to the path
		return home + strings.TrimLeft(path, "~")
	}
	return path
}
