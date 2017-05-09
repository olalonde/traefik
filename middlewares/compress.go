package middlewares

import (
	"github.com/NYTimes/gziphandler"
	"net/http"
)

const (
	headerContentEncoding = "Content-Encoding"
)

// Compress is a middleware that allows redirections
type Compress struct {
}

// ServerHTTP is a function used by negroni
func (c *Compress) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Content-Encoding is already set, skip gzip encoding
	if rw.Header().Get(headerContentEncoding) != "" {
		next.ServeHTTP(rw, r)
		return
	}

	newGzipHandler := gziphandler.GzipHandler(next)
	newGzipHandler.ServeHTTP(rw, r)
}
