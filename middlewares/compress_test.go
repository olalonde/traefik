package middlewares

import (
	// "bytes"
	// "compress/gzip"
	// "fmt"
	// "io"
	"io/ioutil"
	"net"
	"net/http"
	// "net/http/httptest"
	"net/url"
	// "strconv"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/containous/traefik/middlewares"
	"github.com/stretchr/testify/assert"
)

const (
	smallTestBody = "aaabbcaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc"
	testBody      = "aaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbcccaaabbbccc"
)

func TestGzipCompress(t *testing.T) {
	b := []byte(testBody)

	n := negroni.New()

	n.Use(&middlewares.Compress{})

	n.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	}))

	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("failed creating listen socket: %v", err)
	}
	defer ln.Close()
	srv := &http.Server{
		Handler: n,
	}
	go srv.Serve(ln)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/", Scheme: "http", Host: ln.Addr().String()},
		Header: make(http.Header),
		Close:  true,
	}
	req.Header.Set("Accept-Encoding", "gzip")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error making http request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}

	assert.NotEqual(t, testBody, string(body))
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
}

func TestGzipDoubleCompress(t *testing.T) {
	b := []byte(testBody)

	n := negroni.New()

	n.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "already-encoded")
	}))

	n.Use(&middlewares.Compress{})

	n.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	}))

	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("failed creating listen socket: %v", err)
	}
	defer ln.Close()
	srv := &http.Server{
		Handler: n,
	}
	go srv.Serve(ln)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/", Scheme: "http", Host: ln.Addr().String()},
		Header: make(http.Header),
		Close:  true,
	}
	req.Header.Set("Accept-Encoding", "gzip")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error making http request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}

	assert.Equal(t, testBody, string(body))
	assert.Equal(t, "already-encoded", res.Header.Get("Content-Encoding"))
}

/*
func TestGzipDoubleCompress(t *testing.T) {
	b := []byte(testBody)
	compressMiddleware := &middlewares.Compress{}

	handler := compressMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Header().Set("Content-Encoding", "already-encoded")
		w.Write(b)
	}))
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("failed creating listen socket: %v", err)
	}
	defer ln.Close()
	srv := &http.Server{
		Handler: handler,
	}
	go srv.Serve(ln)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/", Scheme: "http", Host: ln.Addr().String()},
		Header: make(http.Header),
		Close:  true,
	}
	req.Header.Set("Accept-Encoding", "gzip")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error making http request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %v", err)
	}

	assert.Equal(t, body, b)
	assert.Equal(t, "already-encoded", res.Header.Get("Content-Encoding"))
}
*/
