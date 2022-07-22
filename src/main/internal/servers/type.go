package servers

import (
	"go-ginApp/src/main/pkg/utils/prof"
	"net/http"
)

var (
	instance *HTTPServer
)

type (
	HTTPServer struct {
		server  *http.Server
		reqCnt  map[string]*prof.CountQPS
		errCnt  map[string]*prof.CountBase
		reqTime map[string]*prof.AverageTimer
	}
)
