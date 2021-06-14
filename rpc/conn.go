package rpc

import (
	"bytes"
	"net/http"
)

// reader wraps byte slices, implementing io.Reader and io.ByteReader.
// Used for self describing data. Meant to never reach EOF.
type reader struct {
	b []byte
	i int64 // read index
}

func newReader(b []byte) *reader {
	return &reader{b: b}
}

// done returns true if underlying data has been completely read.
func (x *reader) done() bool {
	return x.i >= int64(len(x.b))
}

func (x *reader) Read(b []byte) (int, error) {
	n := copy(b, x.b[x.i:])
	x.i += int64(n)

	return n, nil
}

func (x *reader) ReadByte() (byte, error) {
	b := x.b[x.i]
	x.i++

	return b, nil
}

// A clientConn handles byte slice transfer via HTTP POST on the client side.
// Does not implement any thread safety as it's intended to run in WASM (single thread), and Write will always be called before Read.
type clientConn struct {
	url string  // request destination
	r   *reader // holds received data from server; needed for multiple reads
}

func newClientConn(url string) *clientConn {
	return &clientConn{
		url: url,
	}
}

func (x *clientConn) Write(b []byte) (int, error) {
	resp, err := http.Post(x.url, "application/octet-stream", bytes.NewReader(b))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	r := make([]byte, resp.ContentLength)
	resp.Body.Read(r)
	x.r = newReader(r)

	return len(b), nil
}

func (x *clientConn) Read(b []byte) (int, error) {
	return x.r.Read(b)
}

func (x *clientConn) ReadByte() (byte, error) {
	return x.r.ReadByte()
}

func (x *clientConn) Close() error {
	http.DefaultClient.CloseIdleConnections()

	return nil
}

// A serverConn handles byte slice transfer via HTTP POST on the server side.
type serverConn struct {
	mux *http.ServeMux
	srv *http.Server

	r   *reader     // holds requested data; required for reading multiple calls
	rch chan []byte // receives the next incoming data packet
	wch chan []byte // transmits Write data to use as HTTP response

	err error // internal error
}

// newServerConn wraps a new HTTP server listening on the specified port, serving the specified path.
// Provides CORS request handling for the provided origin. Omitted if empty.
func newServerConn(port, path, origin string) *serverConn {
	sc := &serverConn{
		mux: http.NewServeMux(),
		rch: make(chan []byte),
		wch: make(chan []byte),
		r:   newReader(nil), // can immediately use r.done()
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		b := make([]byte, r.ContentLength)
		r.Body.Read(b)
		sc.rch <- b

		w.Write(<-sc.wch)
	})

	if origin != "" {
		handleCORS := func(h http.Handler) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodOptions {
					header := w.Header()
					header.Add("access-control-allow-origin", origin)
					header.Add("access-control-allow-method", http.MethodPost)
					header.Add("access-control-allow-headers", "content-type")

					w.Write([]byte("OK"))
				} else {
					w.Header().Add("access-control-allow-origin", origin)
					h.ServeHTTP(w, r)
				}
			}
		}

		sc.mux.HandleFunc(path, handleCORS(handler))
	} else {
		sc.mux.Handle(path, handler)
	}

	sc.srv = &http.Server{
		Addr:    port,
		Handler: sc.mux,
	}

	return sc
}

func (x *serverConn) Write(b []byte) (int, error) {
	if x.err != nil {
		return 0, x.err
	}

	r := make([]byte, len(b))
	copy(r, b)
	x.wch <- r

	return len(b), nil
}

// sync synchronizes reading and http handling
func (x *serverConn) sync() {
	if x.r.done() {
		x.r = newReader(<-x.rch)
	}
}

func (x *serverConn) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if x.err != nil {
		return 0, x.err
	}

	x.sync()

	return x.r.Read(b)
}

func (x *serverConn) ReadByte() (byte, error) {
	if x.err != nil {
		return 0, x.err
	}

	x.sync()

	return x.r.ReadByte()
}

func (x *serverConn) ListenAndServe() error {
	x.err = x.srv.ListenAndServe()
	return x.err
}

func (x *serverConn) Close() error {
	close(x.rch)
	close(x.wch)

	return x.srv.Close()
}
