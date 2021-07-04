package dom

import (
	"github.com/blitz-frost/wasm/css"
)

type Style struct {
	m map[string]string
}

func MakeStyle(attr ...css.Attribute) Style {
	m := make(map[string]string)
	for _, a := range attr {
		for _, field := range a {
			m[field.Name] = field.Value
		}
	}
	return Style{m}
}

// Fork returns a copy of the target style, with the provided modifications.
func (x Style) Fork(attr ...css.Attribute) Style {
	o := MakeStyle()

	for k, v := range x.m {
		o.m[k] = v
	}

	o.Set(attr...)

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
