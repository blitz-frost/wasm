// Package media wraps the JS MediaDevices API.
package media

import (
	"sync"
	"syscall/js"
	"time"

	"github.com/blitz-frost/bytes"
	"github.com/blitz-frost/io/msg"
	"github.com/blitz-frost/wasm"
)

var (
	media    = wasm.Object(wasm.Global.Get("navigator").Get("mediaDevices"))
	recorder = wasm.ClassGet("MediaRecorder")
	source   = wasm.ClassGet("MediaSource")
)

const (
	VideoInput  DeviceKind = "videoinput"
	AudioInput  DeviceKind = "audioinput"
	AudioOutput DeviceKind = "audiooutput"
)

const (
	User        FacingMode = "user"
	Environment FacingMode = "environment"
	Left        FacingMode = "left"
	Right       FacingMode = "right"
)

const (
	Exact Qualifier = "exact"
	Ideal Qualifier = "ideal"
	Max   Qualifier = "max"
	Min   Qualifier = "min"
)

const (
	None      ResizeMode = "none"
	CropScape ResizeMode = "crop-and-scale"
)

const (
	Audio Kind = "audio"
	Video Kind = "video"
)

type Buffer struct {
	v wasm.Object

	n     int        // js array length
	array wasm.Bytes // copy to JS without repeated allocation
}

func (x *Buffer) Write(b []byte) error {
	if len(b) > x.n {
		x.array = wasm.BytesMake(len(b), len(b))
	}

	slice := x.array.Slice(0, len(b))
	slice.CopyFrom(b)
	x.v.CallRaw("appendBuffer", slice.Value())

	return nil
}

type Device struct {
	Id      string
	GroupId string
}

// Devices returns a slice of all available devices of the specified kind.
//
// Must not be called from the event loop.
func Devices(kind DeviceKind) ([]Device, error) {
	p := wasm.Promise(media.CallRaw("enumerateDevices"))
	allJs, err := p.Await()
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
	v wasm.Object

	active bool
	stop   chan struct{}

	mux sync.Mutex
}

func RecorderMake(s Stream, t Type, audioBitRate, videoBitRate float64) *Recorder {
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

	args := []any{s.V}
	if len(opts) > 0 {
		args = append(args, opts)
	}

	v := recorder.NewRaw(args...)

	x := Recorder{
		v:    v,
		stop: make(chan struct{}),
	}

	return &x
}

// fn - Function({"data": {"arrayBuffer": Function() Promise(ArrayBuffer)}})
func (x *Recorder) DataHandle(fn wasm.Function) {
	x.v.Set("ondataavailable", fn.Value())
}

// DataInterface can be used to create Functions for the DataHandle method.
//
// fn - Function(ArrayBuffer)
func (x *Recorder) DataInterface(fn wasm.Function) wasm.Executor {
	return func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		data := args[0].Get("data")
		arrayPromise := data.Call("arrayBuffer")
		arrayPromise.Call("then", wasm.Value(fn))

		return nil, nil
	}
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

func (x *Recorder) listen(d time.Duration) {
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
	v wasm.Object
}

func settingsMake() Settings {
	v := wasm.ObjectMake()
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
		k := wasm.Object(oJs).Keys()
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
		k := wasm.Object(oJs).Keys()
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
	v wasm.Object

	closeFn wasm.DynamicFunction
	endFn   wasm.DynamicFunction
	openFn  wasm.DynamicFunction
}

func SourceMake() *Source {
	v := source.NewRaw()
	return &Source{
		v: v,
	}
}

func (x *Source) Buffer(t Type) *Buffer {
	s := typeString(t)
	v := x.v.CallRaw("addSourceBuffer", s)
	return &Buffer{
		v: wasm.Object(v),
	}
}

func (x *Source) CloseHandle(fn func()) {
	inter := func(wasm.Value, []wasm.Value) (wasm.Any, error) {
		fn()
		return nil, nil
	}
	x.closeFn.Remake(wasm.Executor(inter))

	x.v.Set("onsourceclose", x.closeFn.Value())
}

func (x *Source) EndHandle(fn func()) {
	inter := func(wasm.Value, []wasm.Value) (wasm.Any, error) {
		fn()
		return nil, nil
	}
	x.endFn.Remake(wasm.Executor(inter))

	x.v.Set("onsourceend", x.endFn.Value())
}

func (x *Source) OpenHandle(fn func()) {
	inter := func(wasm.Value, []wasm.Value) (wasm.Any, error) {
		fn()
		return nil, nil
	}
	x.openFn.Remake(wasm.Executor(inter))

	x.v.Set("onsourceopen", x.openFn.Value())
}

// Url returns a browser URL to the Source object.
func (x *Source) Url() string {
	return wasm.Global.Get("URL").Call("createObjectURL", x.v).String()
}

func (x *Source) Wipe() {
	x.closeFn.Wipe()
	x.endFn.Wipe()
	x.openFn.Wipe()
}

type Stream struct {
	V wasm.Object
}

func (x Stream) Id() string {
	return x.V.Get("id").String()
}

func (x Stream) VideoTracks() []VideoTrack {
	oJs := x.V.CallRaw("getVideoTracks")
	o := make([]VideoTrack, oJs.Length())
	for i := range o {
		o[i].Track = &Track{
			V: wasm.Object(oJs.Index(i)),
		}
	}
	return o
}

type Track struct {
	V wasm.Object

	endFn wasm.DynamicFunction
}

func (x *Track) Enabled() bool {
	return x.V.Get("enabled").Bool()
}

func (x *Track) EnabledSet(v bool) {
	x.V.Set("enabled", v)
}

func (x *Track) Id() string {
	return x.V.Get("id").String()
}

func (x *Track) EndHandle(fn func()) {
	inter := func(wasm.Value, []wasm.Value) (wasm.Any, error) {
		fn()
		return nil, nil
	}
	x.endFn.Remake(wasm.Executor(inter))

	x.V.Set("onended", x.endFn.Value())
}

func (x *Track) Kind() Kind {
	return Kind(x.V.Get("kind").String())
}

func (x *Track) Stop() {
	x.V.CallRaw("stop")
}

func (x *Track) Wipe() {
	x.endFn.Wipe()
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

func VideoSettingsMake() VideoSettings {
	return VideoSettings{settingsMake()}
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

type VideoTrack struct {
	*Track
}

// Must not be called from the event loop.
func (x VideoTrack) Apply(vs VideoSettings) error {
	p := wasm.Promise(x.V.CallRaw("applyConstraints", vs.v))
	_, err := p.Await()
	return err
}

func (x VideoTrack) Capabilities() VideoSettings {
	v := x.V.CallRaw("getCapabilities")
	return VideoSettings{Settings{wasm.Object(v)}}
}

func (x VideoTrack) Settings() VideoSettings {
	v := x.V.CallRaw("getSettings")
	return VideoSettings{Settings{wasm.Object(v)}}
}

type number interface {
	float64 | uint64
}

type single interface {
	bool | string
}

// Get returns a camera video stream.
//
// If a setting is a zero value, it will be ignored. Unmodified settings obtained from a respective make function is equivalent to requesting any stream of that kind.
//
// Must not be called from the event loop.
func Get(video VideoSettings) (Stream, error) {
	con := make(map[string]any)
	if videoVal := wasm.Value(video.v); !videoVal.IsUndefined() {
		k := video.v.Keys()
		if len(k) == 0 {
			con["video"] = true
		} else {
			con["video"] = videoVal
		}
	}

	p := wasm.Promise(media.CallRaw("getUserMedia", con))
	val, err := p.Await()
	return Stream{wasm.Object(val)}, err
}

// GetDisplay returns a display screen stream.
//
// Must not be called from the event loop.
func GetDisplay() (Stream, error) {
	// TODO the call can have an object argument

	p := wasm.Promise(media.CallRaw("getDisplayMedia"))
	v, err := p.Await()
	return Stream{wasm.Object(v)}, err
}

// An ArrayBufferInterface can be used to create Function(ArrayBuffer) that transfer data to a destination.
func ArrayBufferInterfaceMake(dst msg.InputTaker, errorFunc func(error)) wasm.Executor {
	var b []byte

	return func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		buf := wasm.View(args[0])

		n := buf.Len()
		// sometimes we get empty arrays
		if n == 0 {
			return nil, nil
		}
		if len(b) < n {
			b = make([]byte, n)
		}
		b := b[:n]

		buf.CopyTo(b)

		r := (*bytes.Reader)(&b)
		in := msg.Input{
			Read: r.Read,
			Close: func() error {
				return nil
			},
		}
		if err := dst(in); err != nil {
			errorFunc(err)
		}

		return nil, nil
	}
}

func numberGet[T number](x wasm.Object, name string) map[Qualifier]T {
	o := make(map[Qualifier]T)

	oJs := x.Get(name)
	switch oJs.Type() {
	case js.TypeNumber:
		o[Exact] = T(oJs.Float())
	case js.TypeObject:
		k := wasm.Object(oJs).Keys()
		for _, key := range k {
			o[Qualifier(key)] = T(oJs.Get(key).Float())
		}
	}

	return o
}

func numberSet[M ~map[Qualifier]T, T number](x wasm.Object, name string, v M) {
	m := make(map[string]any, len(v))
	for q, a := range v {
		m[string(q)] = a
	}

	x.Set(name, m)
}

func singleSet[T single](x wasm.Object, name string, q Qualifier, v T) {
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
