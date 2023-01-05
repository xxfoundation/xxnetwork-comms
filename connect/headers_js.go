//go:build js && wasm

package connect

import "net/http"

func (wc *webConn) addHeaders(header http.Header) http.Header {
	header.Set("Content-Type", "application/grpc-web+proto")
	header.Add("js.fetch:mode", "no-cors")
	return header
}
