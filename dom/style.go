package dom

import (
	"github.com/blitz-frost/wasm/css"
)

type Style struct {
	m map[string]string
}

func NewStyle(attr ...css.Attribute) Style {
	m := make(map[string]string)
	for _, a := range attr {
		for _, field := range a {
			m[field.Name] = field.Value
		}
	}
	return Style{m}
}

func (x Style) Clone() Style {
	o := NewStyle()

	for k, v := range x.m {
		o.m[k] = v
	}

	return o
}

func (x Style) Include(style Style) {
	for k, v := range style.m {
		x.m[k] = v
	}
}

func (x Style) Set(attr ...css.Attribute) {
	for _, a := range attr {
		for _, field := range a {
			x.m[field.Name] = field.Value
		}
	}
}
