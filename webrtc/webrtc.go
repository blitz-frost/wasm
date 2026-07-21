// Package webrtc wraps the Javascript WebRTC API.
package webrtc

import (
	"github.com/blitz-frost/errors"
	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/media"
)

const (
	BundlePolicyBalanced  BundlePolicy = "balanced"
	BundlePolicyMaxCompat BundlePolicy = "max-compat"
	BundlePolicyMaxBundle BundlePolicy = "max-bundle"
)

type BundlePolicy string

type Channel struct {
	V wasm.Object

	_close     func()
	_closeFn   wasm.DynamicFunction
	_errorFn   wasm.DynamicFunction
	_messageFn wasm.DynamicFunction
	_openFn    wasm.DynamicFunction
}

func channelMake(v wasm.Object) *Channel {
	x := Channel{
		V: v,
	}

	inter := func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		if x._close != nil {
			x._close()
		}

		x._closeFn.Wipe()
		x._errorFn.Wipe()
		x._messageFn.Wipe()
		x._openFn.Wipe()
		return nil, nil
	}
	x._closeFn = wasm.DynamicFunctionMake(inter)

	x.V.Set("onclose", x._closeFn.Value())

	return &x
}

// This begins an asynchronous closing process, and returns immediately.
// A closing and close event will result as part of that process.
// A closing error will trigger an error event, inbetween the closing and close ones.
func (x *Channel) Close() {
	x.V.CallRaw("close")
}

func (x *Channel) CloseHandle(fn func()) {
	x._close = fn
}

func (x *Channel) Label() string {
	v := x.V.Get("label")
	return v.String()
}

func (x *Channel) ErrorHandle(fn func(error)) {
	exec := func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		rtcError := args[0].Get("error")
		detail := rtcError.Get("errorDetail")
		err := errors.Simple(detail.String())
		fn(err)
		return nil, nil
	}
	x._errorFn.Remake(exec)

	x.V.Set("onerror", x._errorFn.Value())
}

func (x *Channel) InputChain(dst func(wasm.Bytes)) {
	exec := func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		// this is an ArrayBuffer by default, unless the channel's binaryType property is changed
		data := args[0].Get("data")
		b := wasm.View(data)
		dst(b)
		return nil, nil
	}
	x._messageFn.Remake(exec)

	x.V.Set("onmessage", x._messageFn.Value())
}

func (x *Channel) OpenHandle(fn func()) {
	exec := func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		fn()
		return nil, nil
	}
	x._openFn.Remake(exec)

	x.V.Set("onopen", x._openFn.Value())
}

func (x *Channel) Write(b wasm.Bytes) error {
	_, err := x.V.Call("send", b.Value())
	return err
}

type ChannelOptions struct {
	v wasm.Object
}

func ChannelOptionsMake() ChannelOptions {
	return ChannelOptions{
		v: wasm.ObjectMake(),
	}
}

func (x ChannelOptions) Id(v uint16) {
	x.v.Set("id", v)
}

func (x ChannelOptions) Negotiated(v bool) {
	x.v.Set("negotiated", v)
}

func (x ChannelOptions) Ordered(v bool) {
	x.v.Set("ordered", v)
}

func (x ChannelOptions) PacketLifetime(ms uint16) {
	x.v.Set("maxPacketLifeTime", ms)
}

func (x ChannelOptions) Protocol(v string) {
	x.v.Set("protocol", v)
}

func (x ChannelOptions) Retransmits(v int) {
	x.v.Set("maxRetransmits", v)
}

// All properties are defined as optional in the JS API, so they may return zero values.
type CodecParameters struct {
	v wasm.Object
}

func (x CodecParameters) Channels() uint {
	v := x.v.Get("channels")
	if v.IsUndefined() {
		return 0
	}
	return uint(v.Int())
}

func (x CodecParameters) ClockRate() uint {
	v := x.v.Get("clockRate")
	if v.IsUndefined() {
		return 0
	}
	return uint(v.Int())
}

func (x CodecParameters) MimeType() string {
	v := x.v.Get("mimeType")
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

func (x CodecParameters) PayloadType() byte {
	v := x.v.Get("payloadType")
	if v.IsUndefined() {
		return 0
	}
	return byte(v.Int())
}

func (x CodecParameters) Sdp() string {
	v := x.v.Get("sdpFmtpLine")
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

type Conn struct {
	V wasm.Object

	_channelFn           wasm.DynamicFunction
	_connectionStateFn   wasm.DynamicFunction
	_iceCandidateFn      wasm.DynamicFunction
	_iceGatheringStateFn wasm.DynamicFunction
	_signalingStateFn    wasm.DynamicFunction
	_trackFn             wasm.DynamicFunction
}

// opt can be the zero value
func ConnMake(opt Config) Conn {
	class := wasm.ClassGet("RTCPeerConnection")

	optVal := wasm.Value(opt.v)

	var args []wasm.Any
	if !optVal.IsUndefined() {
		args = []wasm.Any{optVal}
	}

	return Conn{
		V: class.NewRaw(args...),
	}
}

func (x *Conn) Answer() (SessionDescription, error) {
	r := wasm.Promise(x.V.CallRaw("createAnswer"))
	v, err := r.Await()
	return SessionDescription{wasm.Object(v)}, err
}

// opt can be the zero value
func (x *Conn) Channel(label string, opt ChannelOptions) (*Channel, error) {
	args := []wasm.Any{label}
	if optVal := wasm.Value(opt.v); !optVal.IsUndefined() {
		args = append(args, optVal)
	}

	v, err := x.V.Call("createDataChannel", args...)
	if err != nil {
		return nil, err
	}
	return channelMake(wasm.Object(v)), nil
}

func (x *Conn) ChannelHandle(fn func(*Channel)) {
	inter := func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		v := args[0].Get("channel")
		fn(channelMake(wasm.Object(v)))
		return nil, nil
	}
	x._channelFn.Remake(inter)

	x.V.Set("ondatachannel", x._channelFn.Value())
}

func (x *Conn) Close() {
	x.V.CallRaw("close")

	// order of events is unclear
	// it is only specified that "signalingState" has changed to "closed" by now
	// unclear if that means the respective event was also fired synchronously
	//
	// for now we assume that that is the case, and is also the last event

	x._channelFn.Wipe()
	x._connectionStateFn.Wipe()
	x._iceCandidateFn.Wipe()
	x._iceGatheringStateFn.Wipe()
	x._signalingStateFn.Wipe()
	x._trackFn.Wipe()
}

func (x *Conn) ConnectionStateHandle(fn func(ConnectionState)) {
	x._connectionStateFn.Remake(func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		state := x.V.Get("connectionState").String()
		fn(ConnectionState(state))

		return nil, nil
	})

	x.V.Set("onconnectionstatechange", x._connectionStateFn.Value())
}

func (x *Conn) IceCandidateAdd(candidate IceCandidate) error {
	r := wasm.Promise(x.V.CallRaw("addIceCandidate", wasm.Value(candidate.v)))
	_, err := r.Await()
	return err
}

// Unlike the JS API, this never delivers a null value.
func (x *Conn) IceCandidateHandle(fn func(IceCandidate)) {
	x._iceCandidateFn.Remake(func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		e := wasm.Object(args[0])
		v := e.Get("candidate")
		if v.IsNull() {
			return nil, nil
		}

		fn(IceCandidate{wasm.Object(v)})

		return nil, nil
	})

	x.V.Set("onicecandidate", x._iceCandidateFn.Value())
}

func (x *Conn) IceGatheringState() IceGatheringState {
	v := x.V.Get("iceGatheringState").String()
	return IceGatheringState(v)
}

func (x *Conn) IceGatheringStateHandle(fn func(IceGatheringState)) {
	x._iceGatheringStateFn.Remake(func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		state := x.IceGatheringState()
		fn(state)

		return nil, nil
	})

	x.V.Set("onicegatheringstatechange", x._iceGatheringStateFn.Value())
}

func (x *Conn) LocalDescriptionSet(session SessionDescription) error {
	r := wasm.Promise(x.V.CallRaw("setLocalDescription", wasm.Value(session.v)))
	_, err := r.Await()
	return err
}

func (x *Conn) Offer() (SessionDescription, error) {
	r := wasm.Promise(x.V.CallRaw("createOffer"))
	v, err := r.Await()
	return SessionDescription{wasm.Object(v)}, err
}

func (x *Conn) RemoteDescriptionSet(session SessionDescription) error {
	r := wasm.Promise(x.V.CallRaw("setRemoteDescription", wasm.Value(session.v)))
	_, err := r.Await()
	return err
}

func (x *Conn) SenderRemove(sender Sender) {
	x.V.Call("removeTrack", sender.v)
}

func (x *Conn) SignalingStateHandle(fn func(SignalingState)) {
	x._signalingStateFn.Remake(func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		// the JS event doesn't actually deliver anything

		state := x.V.Get("signalingState").String()
		fn(SignalingState(state))

		return nil, nil
	})

	x.V.Set("onsignalingstatechange", x._signalingStateFn.Value())
}

// opt can be the zero value
func (x *Conn) TrackAdd(track *media.Track, opt media.Stream) (Sender, error) {
	args := []wasm.Any{wasm.Value(track.V)}
	if optVal := wasm.Value(opt.V); !optVal.IsUndefined() {
		args = append(args, optVal)
	}

	v, err := x.V.Call("addTrack", args...)
	return Sender{wasm.Object(v)}, err
}

func (x *Conn) TrackHandle(fn func(*media.Track, []media.Stream, Transceiver)) {
	inter := func(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
		var track media.Track
		track.V = wasm.Object(args[0].Get("track"))

		streamsJs := args[0].Get("streams")
		var streams []media.Stream
		for i, n := 0, streamsJs.Length(); i < n; i++ {
			v := streamsJs.Index(i)
			streams = append(streams, media.Stream{wasm.Object(v)})
		}

		transceiver := Transceiver{
			v: wasm.Object(args[0].Get("transceiver")),
		}

		fn(&track, streams, transceiver)
		return nil, nil
	}
	x._trackFn.Remake(inter)

	x.V.Set("ontrack", x._trackFn.Value())
}

func (x *Conn) TransceiverAddKind(kind media.Kind) (Transceiver, error) {
	v, err := x.V.Call("addTransceiver", string(kind))
	return Transceiver{wasm.Object(v)}, err
}

func (x *Conn) TransceiverAddTrack(track *media.Track) (Transceiver, error) {
	v, err := x.V.Call("addTransceiver", wasm.Value(track.V))
	return Transceiver{wasm.Object(v)}, err
}

const (
	ConnectionStateNew          ConnectionState = "new"
	ConnectionStateConnecting   ConnectionState = "connecting"
	ConnectionStateConnected    ConnectionState = "connected"
	ConnectionStateDisconnected ConnectionState = "disconnected"
	ConnectionStateFailed       ConnectionState = "failed"
	ConnectionStateClosed       ConnectionState = "closed"
)

type ConnectionState string

type Config struct {
	v wasm.Object
}

func ConfigMake() Config {
	return Config{
		v: wasm.ObjectMake(),
	}
}

func (x Config) BundlePolicy(v BundlePolicy) {
	x.v.Set("bundlePolicy", string(v))
}

func (x Config) IceServers(v []IceServer) {
	a := wasm.ArrayMake(len(v))
	for i, server := range v {
		a.IndexSet(i, wasm.Value(server.v))
	}
	x.v.Set("iceServers", wasm.Value(a))
}

func (x Config) IceTransportPolicy(v IceTransportPolicy) {
	x.v.Set("iceTransportPolicy", string(v))
}

const (
	DirectionBoth    Direction = "sendrecv"
	DirectionNone    Direction = "inactive"
	DirectionReceive Direction = "recvonly"
	DirectionSend    Direction = "sendonly"
	DirectionStopped Direction = "stopped" // this must not be manually set on a Transceiver
)

type Direction string

type EncodingParameters struct {
	v wasm.Object
}

func (x EncodingParameters) Active() bool {
	return x.v.Get("active").Bool()
}

func (x EncodingParameters) BitrateMax() uint {
	v := x.v.Get("maxBitrate")
	return uint(v.Int())
}

func (x EncodingParameters) BitrateMaxSet(br uint) {
	x.v.Set("maxBitrate", br)
}

func (x EncodingParameters) Downscale() float64 {
	return x.v.Get("scaleResolutionDownBy").Float()
}

// Only for video tracks.
// factor must be >= 1 and is applied to both image dimensions.
func (x EncodingParameters) DownscaleSet(factor float64) {
	x.v.Set("scaleResolutionDownBy", factor)
}

func (x EncodingParameters) FramerateMax() float64 {
	return x.v.Get("maxFramerate").Float()
}

func (x EncodingParameters) FramerateMaxSet(fps float64) {
	x.v.Set("maxFramerate", fps)
}

func (x EncodingParameters) PtimeSet(ms uint) {
	x.v.Set("ptime", ms)
}

type IceCandidate struct {
	v wasm.Object
}

func IceCandidateMake(info IceCandidateInfo) (IceCandidate, error) {
	class := wasm.ClassGet("RTCIceCandidate")
	v, err := class.New(wasm.Value(info.v))
	return IceCandidate{v}, err
}

func (x IceCandidate) Info() IceCandidateInfo {
	v := x.v.CallRaw("toJSON")
	return IceCandidateInfo{wasm.Object(v)}
}

// All of its properties are optional.
type IceCandidateInfo struct {
	v wasm.Object
}

func IceCandidateInfoMake() IceCandidateInfo {
	return IceCandidateInfo{
		v: wasm.ObjectMake(),
	}
}

func (x IceCandidateInfo) Candidate() string {
	v := x.v.Get("candidate")
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

func (x IceCandidateInfo) CandidateSet(v string) {
	x.v.Set("candidate", v)
}

// Unlike the JS API, uses "" instead of null.
// The intent for null was possibly to allow the empty string as valid, but we assume that it is in fact not valid.
func (x IceCandidateInfo) MediaId() string {
	v := x.v.Get("sdpMid")
	if v.IsUndefined() || v.IsNull() {
		return ""
	}
	return v.String()
}

func (x IceCandidateInfo) MediaIdSet(v string) {
	var vv wasm.Any
	if v != "" {
		vv = v
	}
	x.v.Set("sdpMid", vv)
}

// Unlike the JS API, uses -1 instead of null.
// The index is supposed to be a positive integer, so -1 is clearly invalid.
func (x IceCandidateInfo) MediaLineIndex() int {
	v := x.v.Get("sdpMLineIndex")
	if v.IsUndefined() || v.IsNull() {
		return -1
	}
	return v.Int()
}

func (x IceCandidateInfo) MediaLineIndexSet(v int) {
	var vv wasm.Any
	if v >= 0 {
		vv = v
	}
	x.v.Set("sdpMLineIndex", vv)
}

// Returns "" when not present, which is guaranteed to be invalid (since a valid value must contain at least 3 random bytes).
// This one has no null case; it must not be set when absent.
func (x IceCandidateInfo) UsernameFragment() string {
	v := x.v.Get("usernameFragment")
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

func (x IceCandidateInfo) UsernameFragmentSet(v string) {
	x.v.Set("usernameFragment", v)
}

const (
	IceGatheringStateNew       IceGatheringState = "new"
	IceGatheringStateGathering IceGatheringState = "gathering"
	IceGatheringStateComplete  IceGatheringState = "complete"
)

type IceGatheringState string

type IceServer struct {
	v wasm.Object
}

func IceServerMake(urls ...string) IceServer {
	var urlsVal wasm.Any
	if n := len(urls); n == 1 {
		urlsVal = urls[0]
	} else {
		a := wasm.ArrayMake(n)
		for i, url := range urls {
			a.IndexSet(i, url)
		}
		urlsVal = wasm.Value(a)
	}

	v := wasm.ObjectMake()
	v.Set("urls", urlsVal)

	return IceServer{v}
}

func (x IceServer) Credential(v string) {
	x.v.Set("credential", v)
}

func (x IceServer) Username(v string) {
	x.v.Set("username", v)
}

const (
	IceTransportPolicyAll    IceTransportPolicy = "all"
	IceTransportPolicyPublic IceTransportPolicy = "public"
	IceTransportPolicyRelay  IceTransportPolicy = "relay"
)

type IceTransportPolicy string

type SendParameters struct {
	v wasm.Object
}

func (x SendParameters) Codecs() []CodecParameters {
	codecs := x.v.Get("codecs")

	n := codecs.Length()
	o := make([]CodecParameters, n)
	for i := range n {
		v := codecs.Index(i)
		o[i] = CodecParameters{wasm.Object(v)}
	}

	return o
}

// Modify the return values directly, then call Sender.ParametersSet(x).
func (x SendParameters) Encodings() []EncodingParameters {
	encodings := x.v.Get("encodings")

	n := encodings.Length()
	o := make([]EncodingParameters, n)
	for i := range n {
		v := encodings.Index(i)
		o[i] = EncodingParameters{wasm.Object(v)}
	}

	return o
}

type Sender struct {
	v wasm.Object
}

func (x Sender) Parameters() SendParameters {
	v := x.v.CallRaw("getParameters")
	return SendParameters{wasm.Object(v)}
}

// Must be called with the return value of the last Parameters method call.
func (x Sender) ParametersSet(params SendParameters) error {
	promise := wasm.Promise(x.v.CallRaw("setParameters", wasm.Value(params.v)))
	_, err := promise.Await()
	return err
}

func (x Sender) TrackReplace(track *media.Track) error {
	promise := wasm.Promise(x.v.CallRaw("replaceTrack", track.V))
	_, err := promise.Await()
	return err
}

type SessionDescription struct {
	v wasm.Object
}

func SessionDescriptionMake(typ SessionDescriptionType, protocol string) SessionDescription {
	v := wasm.ObjectMake()
	v.Set("type", string(typ))
	v.Set("sdp", protocol)
	return SessionDescription{v}
}

func (x SessionDescription) Protocol() string {
	return x.v.Get("sdp").String()
}

func (x SessionDescription) Type() SessionDescriptionType {
	v := x.v.Get("type").String()
	return SessionDescriptionType(v)
}

const (
	SessionDescriptionTypeAnswer      SessionDescriptionType = "answer"
	SessionDescriptionTypeOffer       SessionDescriptionType = "offer"
	SessionDescriptionTypeProvisional SessionDescriptionType = "pranswer"
	SessionDescriptionTypeRollback    SessionDescriptionType = "rollback"
)

type SessionDescriptionType string

const (
	SignalingStateStable       SignalingState = "stable"
	SignalingStateLocalOffer   SignalingState = "have-local-offer"
	SignalingStateRemoteOffer  SignalingState = "have-remote-offer"
	SignalingStateLocalAnswer  SignalingState = "have-local-pranswer"
	SignalingStateRemoteAnswer SignalingState = "have-remote-pranswer"
	SignalingStateClosed       SignalingState = "closed"
)

type SignalingState string

type Transceiver struct {
	v wasm.Object
}

// Direction returns the actual current direction, which might differ from the previously set one.
func (x Transceiver) Direction() Direction {
	o := x.v.Get("currentDirection").String()
	return Direction(o)
}

// DirectionSet sets the desired direction, which is not applied immediately.
// Must not be called with [DirectionStopped], use [Transceiver.Stop] instead.
func (x Transceiver) DirectionSet(dir Direction) {
	x.v.Set("direction", string(dir))
}

func (x Transceiver) Stop() {
	x.v.CallRaw("stop")
}
