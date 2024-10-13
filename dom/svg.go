package dom

import (
	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/svg"
)

type Svg struct {
	Element
}

func SvgMake() Svg {
	return Svg{Element{doc.CallRaw("createElementNS", "http://www.w3.org/2000/svg", "svg")}}
}

func (x Svg) Append(e ...svg.Element) {
	for _, elem := range e {
		x.Call("appendChild", wasm.Value(elem))
	}
}

func (x Svg) Sub(i int) wasm.Value {
	return x.Get("children").Index(i)
}
