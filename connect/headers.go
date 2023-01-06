//go:build !js || !wasm

package connect

import "net/http"

func (wc *webConn) addHeaders(header http.Header) http.Header {
	return header
}
