package connect

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"io"
	"net/http"
)

type httpTraceWrapper struct {
	ws *grpcweb.WrappedGrpcServer
}

func (htw *httpTraceWrapper) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodTrace {
		// TODO or do we just return a 200?  writing the body back is proper http trace protocol but how much do we actually care?
		received, err := io.ReadAll(req.Body)
		if err != nil {
			respWriter.WriteHeader(500)
			return
		}
		_, err = respWriter.Write(received)
		if err != nil {
			respWriter.WriteHeader(500)
			return
		}
		return
	} else {
		htw.ws.ServeHTTP(respWriter, req)
	}
}
