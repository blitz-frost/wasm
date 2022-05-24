// Package media wraps the JS MediaDevices API.
package media

import (
	"errors"
	"sync"
	"syscall/js"
	"time"

	"github.com/blitz-frost/io"
	"github.com/blitz-frost/wasm"
)

var (
	media    = js.Global().Get("navigator").Get("mediaDevices")
	recorder = js.Global().Get("MediaRecorder")
	source   = js.Global().Get("MediaSource")
)

const (
	VideoInput  DeviceKind = "videoinput"
	AudioInput             = "audioinput"
	AudioOutput            = "audiooutput"
)

const (
	User        FacingMode = "user"
	Environment            = "environment"
	Left                   = "left"
	Right                  = "right"
)

const (
	Exact Qualifier = "exact"
	Ideal           = "ideal"
	Max             = "max"
	Min             = "min"
)

const (
	None      ResizeMode = "none"
	CropScape            = "crop-and-scale"
)

const (
	Audio Kind = "audio"
	Video      = "video"
)

type Buffer struct {
	v js.Value

	n     int        // js array length
	array wasm.Bytes // copy to JS without repeated allocation
}

func newBuffer(v js.Value) *Buffer {
	return &Buffer{
		v: v,
	}
}

func (x *Buffer) Write(b []byte) error {
	if len(b) > x.n {
		x.array = wasm.MakeBytes(len(b))
	}

	slice := x.array.Slice(0, len(b))
	slice.CopyFrom(b)
	x.v.Call("appendBuffer", slice.Js())

	return nil
}

type Device struct {
	Id      string
	GroupId string
}

// Devices returns a slice of all available devices of the specified kind.
func Devices(kind DeviceKind) ([]Device, error) {
	allJs, err := wasm.Await(media.Call("enumerateDevices"))
	if err != nil {
		return nil, err
	}

	var o []Device
	for i, n := 0, allJs.Length(); i < n; i++ {
		deviceJs := allJs.Index(i)
		if deviceJs.Get("kind").String() == string(kind) {
			o = append(o, Device{
				Id:      deviceJs.Get("deviceId").String(),
				GroupId: deviceJs.Get("groupId").String(),
			})
		}
	}

	return o, nil
}

type DeviceKind string

type FacingMode string

type Float map[Qualifier]float64

type Kind string

type Qualifier string

type Recorder struct {
	v js.Value

	onArray   js.Func // should be more efficient than awaiting the onData promise
	onData    js.Func
	onErrorJs js.Func // onerror event listener

	onError func(error) // also used for dst.Write errors

	dst io.Writer
	buf []byte // receive recorded bytes without repeated allocation

	active bool
	stop   chan struct{}

	mux sync.Mutex
}

func NewRecorder(s Stream, t Type, audioBitRate, videoBitRate float64) *Recorder {
	// options
	opts := make(map[string]any)
	if t != nil {
		opts["mimeType"] = typeString(t)
	}
	if audioBitRate != 0 {
		opts["audioBitsPerSecond"] = audioBitRate
	}
	if videoBitRate != 0 {
		opts["videoBitsPerSecond"] = videoBitRate
	}

	args := []any{s.v}
	if len(opts) > 0 {
		args = append(args, opts)
	}

	v := recorder.New(args...)

	x := Recorder{
		v:       v,
		onError: func(error) {},
		dst:     io.VoidWriter{},
		stop:    make(chan struct{}),
	}

	x.onErrorJs = js.FuncOf(func(this js.Value, args []js.Value) any {
		errJs := args[0].Get("error")
		msg := errJs.Get("message").String()
		err := errors.New(msg)
		x.onError(err)

		return nil
	})
	x.onArray = js.FuncOf(func(this js.Value, args []js.Value) any {
		buf := wasm.View(args[0])

		n := buf.Length()
		// sometimes we get empty arrays
		if n == 0 {
			return nil
		}
		if len(x.buf) < n {
			x.buf = make([]byte, n)
		}
		b := x.buf[:n]

		buf.CopyTo(b)
		if err := x.dst.Write(b); err != nil {
			x.onError(err)
		}

		return nil
	})
	x.onData = js.FuncOf(func(this js.Value, args []js.Value) any {
		data := args[0].Get("data")
		arrayPromise := data.Call("arrayBuffer")
		arrayPromise.Call("then", x.onArray)

		return nil
	})

	v.Set("ondataavailable", x.onData)

	return &x
}

func (x *Recorder) Chain(w io.Writer) {
	x.dst = w
}

func (x Recorder) ChainGet() io.Writer {
	return x.dst
}

func (x *Recorder) OnError(fn func(error)) {
	x.onError = fn
}

func (x *Recorder) Pause() {
	x.mux.Lock()
	defer x.mux.Unlock()

	if !x.active {
		return
	}
	x.active = false
	x.stop <- struct{}{}

	x.v.Call("pause")
}

func (x Recorder) Release() {
	x.onArray.Release()
	x.onData.Release()
	x.onErrorJs.Release()
}

func (x *Recorder) Resume(d time.Duration) {
	x.mux.Lock()
	defer x.mux.Unlock()

	if x.active {
		return
	}
	x.active = true

	x.v.Call("resume")

	go x.listen(d)
}

// Start starts recording. Writes output every d.
func (x *Recorder) Start(d time.Duration) {
	x.mux.Lock()
	defer x.mux.Unlock()

	if x.active {
		return
	}
	x.active = true

	x.v.Call("start")

	go x.listen(d)
}

func (x *Recorder) Stop() {
	x.mux.Lock()
	defer x.mux.Unlock()

	if !x.active {
		return
	}
	x.active = false
	x.stop <- struct{}{}

	x.v.Call("stop")
}

func (x Recorder) listen(d time.Duration) {
	t := time.NewTicker(d)
	for {
		select {
		case <-x.stop:
			t.Stop()
			return
		case <-t.C:
			x.v.Call("requestData")
		}
	}
}

type ResizeMode string

// Settings defines a set of properties common to all stream types.
type Settings struct {
	v js.Value
}

func makeSettings() Settings {
	v := js.ValueOf(map[string]any{})
	return Settings{v}
}

func (x Settings) Device() (Qualifier, string) {
	return x.stringGet("deviceId")
}

func (x Settings) DeviceSet(q Qualifier, id string) {
	x.stringSet("deviceId", q, id)
}

func (x Settings) Group() (Qualifier, string) {
	return x.stringGet("groupId")
}

func (x Settings) GroupSet(q Qualifier, id string) {
	x.stringSet("groupId", q, id)
}

func (x Settings) boolGet(name string) (Qualifier, bool) {
	oJs := x.v.Get(name)
	switch oJs.Type() {
	case js.TypeBoolean:
		return Exact, oJs.Bool()
	case js.TypeObject:
		k := wasm.Keys(oJs)
		return Qualifier(k[0]), oJs.Get(k[0]).Bool()
	}

	return "", false
}

func (x Settings) boolSet(name string, q Qualifier, v bool) {
	singleSet(x.v, name, q, v)
}

func (x Settings) floatGet(name string) Float {
	return Float(numberGet[float64](x.v, name))
}

func (x Settings) floatSet(name string, v Float) {
	numberSet(x.v, name, v)
}

func (x Settings) stringGet(name string) (Qualifier, string) {
	oJs := x.v.Get(name)
	switch oJs.Type() {
	case js.TypeString:
		return Exact, oJs.String()
	case js.TypeObject:
		k := wasm.Keys(oJs)
		return Qualifier(k[0]), oJs.Get(k[0]).String()
	}

	return "", ""
}

func (x Settings) stringSet(name string, q Qualifier, v string) {
	singleSet(x.v, name, q, v)
}

func (x Settings) uintGet(name string) Uint {
	return Uint(numberGet[uint64](x.v, name))
}

func (x Settings) uintSet(name string, v Uint) {
	numberSet(x.v, name, v)
}

type Source struct {
	v js.Value

	onClose js.Func
	onEnd   js.Func
	onOpen  js.Func
}

func NewSource() *Source {
	v := source.New()

	return &Source{
		v: v,
	}
}

func (x Source) NewBuffer(t Type) *Buffer {
	s := typeString(t)
	v := x.v.Call("addSourceBuffer", s)
	return newBuffer(v)
}

func (x *Source) OnClose(fn func()) {
	x.onClose.Release()
	x.onClose = js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	})
	x.v.Set("onsourceclose", x.onClose)
}

func (x *Source) OnEnd(fn func()) {
	x.onEnd.Release()
	x.onEnd = js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	})
	x.v.Set("onsourceended", x.onEnd)
}

func (x *Source) OnOpen(fn func()) {
	x.onOpen.Release()
	x.onOpen = js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	})
	x.v.Set("onsourceopen", x.onOpen)
}

func (x Source) Release() {
	x.onClose.Release()
	x.onEnd.Release()
	x.onOpen.Release()
}

// Url returns a browser URL to the Source object.
func (x Source) Url() string {
	return js.Global().Get("URL").Call("createObjectURL", x.v).String()
}

type Stream struct {
	v js.Value
}

func AsStream(v js.Value) Stream {
	return Stream{v}
}

func (x Stream) Js() js.Value {
	return x.v
}

func (x Stream) VideoTracks() []VideoTrack {
	oJs := x.v.Call("getVideoTracks")
	o := make([]VideoTrack, oJs.Length())
	for i := range o {
		o[i] = VideoTrack{oJs.Index(i)}
	}
	return o
}

type Track struct {
	v js.Value
}

func AsTrack(v js.Value) Track {
	return Track{v}
}

func (x Track) Js() js.Value {
	return x.v
}

type Type interface {
	Kind() Kind
	Format() string
	Codec() string
}

type Uint map[Qualifier]uint64

type VideoSettings struct {
	Settings
}

func MakeVideoSettings() VideoSettings {
	return VideoSettings{makeSettings()}
}

func (x VideoSettings) AspectRatio() Float {
	return x.floatGet("aspectRatio")
}

func (x VideoSettings) AspectRatioSet(f Float) {
	x.floatSet("aspectRatio", f)
}

func (x VideoSettings) FacingMode() (Qualifier, FacingMode) {
	q, o := x.stringGet("facingMode")
	return q, FacingMode(o)
}

func (x VideoSettings) FacingModeSet(q Qualifier, fm FacingMode) {
	x.stringSet("facingMode", q, string(fm))
}

func (x VideoSettings) FrameRate() Float {
	return x.floatGet("frameRate")
}

func (x VideoSettings) FrameRateSet(f Float) {
	x.floatSet("frameRate", f)
}

func (x VideoSettings) Height() Uint {
	return x.uintGet("height")
}

func (x VideoSettings) HeightSet(u Uint) {
	x.uintSet("height", u)
}

func (x VideoSettings) ResizeMode() ResizeMode {
	// unlike other constraints, resizeMode can't have a qualifier

	s := x.v.Get("resizeMode").String()
	return ResizeMode(s)
}

func (x VideoSettings) ResizeModeSet(rm ResizeMode) {
	x.v.Set("resizeMode", string(rm))
}

func (x VideoSettings) Width() Uint {
	return x.uintGet("width")
}

func (x VideoSettings) WidthSet(u Uint) {
	x.uintSet("width", u)
}

type VideoTrack Track

func (x VideoTrack) Apply(vs VideoSettings) error {
	_, err := wasm.Await(x.v.Call("applyConstraints", vs.v))
	return err
}

func (x VideoTrack) Capabilities() VideoSettings {
	v := x.v.Call("getCapabilities")
	return VideoSettings{Settings{v}}
}

func (x VideoTrack) Settings() VideoSettings {
	v := x.v.Call("getSettings")
	return VideoSettings{Settings{v}}
}

type number interface {
	float64 | uint64
}

type single interface {
	bool | string
}

// If a setting is a zero value, it will be ignored. Unmodified settings obtained from a respective make function is equivalent to requesting any stream of that kind.
func Get(video VideoSettings) (Stream, error) {
	con := make(map[string]any)
	if !video.v.IsUndefined() {
		k := wasm.Keys(video.v)
		if len(k) == 0 {
			con["video"] = true
		} else {
			con["video"] = video.v
		}
	}

	val, err := wasm.Await(media.Call("getUserMedia", con))
	return Stream{val}, err
}

func numberGet[T number](x js.Value, name string) map[Qualifier]T {
	o := make(map[Qualifier]T)

	oJs := x.Get(name)
	switch oJs.Type() {
	case js.TypeNumber:
		o[Exact] = T(oJs.Float())
	case js.TypeObject:
		k := wasm.Keys(oJs)
		for _, key := range k {
			o[Qualifier(key)] = T(oJs.Get(key).Float())
		}
	}

	return o
}

func numberSet[M ~map[Qualifier]T, T number](x js.Value, name string, v M) {
	m := make(map[string]any, len(v))
	for q, a := range v {
		m[string(q)] = a
	}

	x.Set(name, m)
}

func singleSet[T single](x js.Value, name string, q Qualifier, v T) {
	m := map[string]any{
		string(q): v,
	}

	x.Set(name, m)
}

func typeString(t Type) string {
	o := string(t.Kind())
	o += "/" + t.Format()
	c := t.Codec()
	if c != "" {
		o += "; codecs=\"" + c + "\""
	}

	return o
}
