// Package webm implements the WebM media type.
package webm

import (
	"strconv"

	"github.com/blitz-frost/wasm/media"
)

const (
	Opus   Audio = "opus"
	Vorbis Audio = "vorbis"
)

const (
	Depth8  BitDepth = "08"
	Depth10 BitDepth = "10"
	Depth12 BitDepth = "12"
)

const (
	Profile8Bit420 Profile = "00"
	Profile8Bit    Profile = "01"
	ProfileAny420  Profile = "02" // 8, 10, 12 bit depth
	ProfuleAny     Profile = "03"
)

type Audio string

func (x Audio) Kind() media.Kind {
	return media.Audio
}

func (x Audio) Format() string {
	return "webm"
}

func (x Audio) Codec() string {
	return string(x)
}

type BitDepth string

// Level 3 is {3, 0}
// Level 6.1 is {6, 1}
type Level [2]byte

type Profile string

// A Video represents a webm video media type. The zero value is invalid.
type Video struct {
	codec   string
	profile string
	level   string
	depth   string
	audio   string
}

func VP8() Video {
	return Video{
		codec: "vp8",
	}
}

func VP9() Video {
	return Video{
		codec: "vp9",
	}
}

// Audio specifies an audio codec to use for AV streams.
func (x *Video) Audio(a Audio) {
	x.audio = string(a)
}

// Set specifies the 3 required VP parameters.
// A Video value is valid for use even without calling this method.
func (x *Video) Set(p Profile, l Level, d BitDepth) {
	x.profile = string(p)
	x.level = strconv.Itoa(int(l[0])) + strconv.Itoa(int(l[1]))
	x.depth = string(d)
}

func (x Video) Kind() media.Kind {
	return media.Video
}

func (x Video) Format() string {
	return "webm"
}

func (x Video) Codec() string {
	o := x.codec

	if x.profile != "" {
		// if profile is set, then level and depth should also be
		o += "." + x.profile + "." + x.level + "." + x.depth
	}

	if x.audio != "" {
		o += "," + x.audio
	}

	return o
}
