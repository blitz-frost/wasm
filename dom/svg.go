package dom

import (
	"syscall/js"

	"github.com/blitz-frost/wasm/svg"
)

type Svg struct {
	Element
}

func SvgMake() Svg {
	return Svg{Element{doc.Call("createElementNS", "http://www.w3.org/2000/svg", "svg")}}
}

func (x Svg) Append(e ...svg.Element) {
	for _, elem := range e {
		x.Call("appendChild", elem.JSValue())
	}
}

func (x Svg) Sub(i int) js.Value {
	return x.Get("children").Index(i)
}
