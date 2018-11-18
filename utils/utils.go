package utils

import (
	"github.com/mitchellh/go-homedir"
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"strings"
)

// Given a path, replace a "~" character
// with the home directory to return a full file path
func GetFullPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}
		// Append the home directory to the path
		return home + strings.TrimLeft(path, "~")
	}
	return path
}
