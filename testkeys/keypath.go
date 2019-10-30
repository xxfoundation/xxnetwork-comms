package testkeys

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/utils"
	"path/filepath"
	"runtime"
)

func getDirForFile() string {
	// Get the filename we're in
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}

// These functions are used to cover TLS connection code in tests
func GetNodeCertPath() string {
	return filepath.Join(getDirForFile(), "cmix.rip.crt")
}

func GetNodeKeyPath() string {
	return filepath.Join(getDirForFile(), "cmix.rip.key")
}

func GetGatewayCertPath() string {
	return filepath.Join(getDirForFile(), "gateway.cmix.rip.crt")
}

func GetGatewayKeyPath() string {
	return filepath.Join(getDirForFile(), "gateway.cmix.rip.key")
}

func LoadFromPath(path string) []byte {
	data, err := utils.ReadFile(path)
	if err != nil {
		jww.FATAL.Panicf("***Check your key!***\nFailed to read file at %s: %+v", path, err)
	}
	return data
}
