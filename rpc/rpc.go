// Package rpc provides bridging between two Go programs using HTTP POST. Uses encoding/gob for encoding.
package rpc

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/blitz-frost/wasm/wire"
)

// A bit of a shenanigan in order to always have the error interface type at hand.
// Used to check if types implement error.
var errorType reflect.Type = reflect.TypeOf(new(error)).Elem()

// A procedure wraps functions to be usable by the Server type.
// Underlying function may have a final error output. All its other return values and arguments must be concrete types.
type procedure struct {
	f        reflect.Value  // underlying function
	inType   []reflect.Type // input types
	numOut   int            // output number, excluding error
	hasError bool           // true if function returns an error
}

// newProcedure fails if f has non-concrete inputs or outputs, excluding a potential final error output.
func newProcedure(f interface{}) (*procedure, error) {
	t := reflect.TypeOf(f)

	if t.Kind() != reflect.Func {
		return nil, errors.New("not a function")
	}

	proc := &procedure{
		f: reflect.ValueOf(f),
	}

	// check inputs
	numIn := t.NumIn()
	proc.inType = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		inType := t.In(i)
		if inType.Kind() == reflect.Interface {
			return nil, errors.New("function has interface argument")
		}
		proc.inType[i] = inType
	}

	// check error
	numOut := t.NumOut()
	if numOut > 0 {
		if t.Out(numOut - 1).Implements(errorType) {
			proc.hasError = true
			numOut--
		}
	}
	proc.numOut = numOut
	// check rest of outputs
	for i := 0; i < numOut; i++ {
		if t.Out(i).Kind() == reflect.Interface {
			return nil, errors.New("function has interface return value")
		}
	}

	return proc, nil
}

// call executes the underlying function with the provided input.
func (x *procedure) call(in []reflect.Value) ([]reflect.Value, error) {
	r := x.f.Call(in)
	if x.hasError && !r[len(r)-1].IsNil() {
		return nil, r[len(r)-1].Interface().(error)
	}

	return r[:x.numOut], nil
}

// buffer is a simple staging Writer implementation.
// Used because gob.Encoder uses multiple Write calls when encoding, while the connection expects a single Write.
type buffer []byte

func (x *buffer) Write(b []byte) (int, error) {
	*x = append(*x, b...)
	return len(b), nil
}

func (x *buffer) WriteTo(w io.Writer) error {
	_, err := w.Write(*x)
	*x = (*x)[:0]
	return err
}

// Client represents an RPC Client.
// It is not concurrent safe.
type Client struct {
	buf  *buffer // used to buffer the encoding of whole values before sending
	conn *clientConn
	enc  *wire.Encoder
	dec  *wire.Decoder
}

func NewClient(url string) *Client {
	cli := &Client{
		buf:  new(buffer),
		conn: newClientConn(url),
	}
	cli.enc = wire.NewEncoder(cli.buf)
	cli.dec = wire.NewDecoder(cli.conn)

	return cli
}

// Bind generates an RPC function through the Client, and stores it inside the given function pointer.
//
// Calls for this function will be made under the given name.
//
// fptr must be a non nil pointer to a function that returns an error as a final return value.
// All other returns values and arguments must be concrete types.
func (x *Client) Bind(name string, fptr interface{}) error {
	fv := reflect.ValueOf(fptr).Elem()
	ft := fv.Type()

	// check inputs
	numIn := ft.NumIn()
	for i := 0; i < numIn; i++ {
		if ft.In(i).Kind() == reflect.Interface {
			return errors.New("function has interface argument")
		}
	}

	// check error
	numOut := ft.NumOut() - 1
	if numOut < 0 || !ft.Out(numOut).Implements(errorType) {
		return errors.New("function does not return an error")
	}

	// check rest of outputs
	outTypes := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		outType := ft.Out(i)
		if outType.Kind() == reflect.Interface {
			return errors.New("function has interface output")
		}
		outTypes[i] = outType
	}

	fn := func(args []reflect.Value) (results []reflect.Value) {
		// prepare pointers to return data types, except the error
		var err error
		results = make([]reflect.Value, numOut+1)
		for i := 0; i < numOut; i++ {
			results[i] = reflect.New(outTypes[i])
		}
		// defer return data pointers before returning
		defer func() {
			for i := 0; i < numOut; i++ {
				results[i] = results[i].Elem()
			}
			if err != nil {
				results[numOut] = reflect.ValueOf(err)
			} else {
				results[numOut] = reflect.Zero(errorType)
			}
		}()

		// encode name and values
		if err = x.enc.Encode(name); err != nil {
			return
		}
		for _, v := range args {
			if err = x.enc.EncodeValue(v); err != nil {
				return
			}
		}

		// send call
		if err = x.buf.WriteTo(x.conn); err != nil {
			return
		}

		// check return error
		var errStr string
		if err = x.dec.Decode(&errStr); err != nil {
			return
		}
		if errStr != "" {
			err = errors.New(errStr)
			return
		}

		// decode return values
		for i := 0; i < numOut; i++ {
			if err = x.dec.DecodeValue(results[i]); err != nil {
				return
			}
		}

		return
	}

	fv.Set(reflect.MakeFunc(ft, fn))
	return nil
}

// Server represents an RPC Server.
type Server struct {
	buf  *buffer
	conn *serverConn
	enc  *wire.Encoder
	dec  *wire.Decoder

	procs map[string]*procedure
}

// NewServer returns a new Server that will listen on the specified port and serve the specified path.
// Will set up CORS handling for the specified origin. origin may be empty, in which case CORS is omitted.
func NewServer(port, path, origin string) *Server {
	srv := &Server{
		buf:   new(buffer),
		conn:  newServerConn(port, path, origin),
		procs: make(map[string]*procedure),
	}
	srv.enc = wire.NewEncoder(srv.buf)
	srv.dec = wire.NewDecoder(srv.conn)

	return srv
}

// Register makes f available as a procedure on the server, under the provided name.
//
// f may have a final error output, which will be propagated in a more streamlined fashion to clients.
// Any input or output interface types, except for the error, must have their concrete types registered with RegisterType before use.
//
// Only the return values are communicated to clients, even if the arguments are mutable.
func (x *Server) Register(name string, f interface{}) error {
	p, err := newProcedure(f)
	if err != nil {
		return err
	}

	x.procs[name] = p

	return nil
}

// Serve the next incoming remote call.
func (x *Server) Serve() error {
	var name string
	if err := x.dec.Decode(&name); err != nil {
		return errors.New("name decode error: " + err.Error())
	}

	p, ok := x.procs[name]
	if !ok {
		return errors.New("invalid name")
	}

	// decode expected arguments
	in := make([]reflect.Value, len(p.inType))
	for i := 0; i < len(p.inType); i++ {
		v := reflect.New(p.inType[i])
		if err := x.dec.DecodeValue(v); err != nil {
			return errors.New("input decode error: " + err.Error())
		}
		in[i] = v.Elem()
	}

	out, err := p.call(in)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	// encode error
	if err := x.enc.Encode(errStr); err != nil {
		return errors.New("error encode error: " + err.Error())
	}

	// only encode rest if error is nil
	if err == nil {
		for i := 0; i < len(out); i++ {
			if err := x.enc.EncodeValue(out[i]); err != nil {
				return errors.New("return value encode error: " + err.Error())
			}
		}
	}

	err = x.buf.WriteTo(x.conn)
	if err != nil {
		return errors.New("flush error: " + err.Error())
	}

	return nil
}

// ListenAndServe starts server operation. Blocks until there is an error.
func (x *Server) ListenAndServe() error {
	var err error

	go func() {
		for err == nil {
			if serveErr := x.Serve(); serveErr != nil {
				fmt.Println("rpc serve error: ", serveErr)
			}
		}
	}()

	err = x.conn.ListenAndServe()
	return err
}
