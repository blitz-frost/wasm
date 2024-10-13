package dom

import (
	"errors"
	"net/url"

	"github.com/blitz-frost/wasm"
)

var (
	doc      = wasm.Object(wasm.Global.Get("document"))
	location = wasm.Object(wasm.Global.Get("location"))
)

// ElementById returns the element with the given ID in the document.
// Returns an error if the ID doesn't exist.
func ElementById(id string) (Element, error) {
	elem := doc.CallRaw("getElementById", id)
	if elem.IsNull() {
		return Element{}, errors.New(id + " not found")
	}
	return Element{elem}, nil
}

// ElementsByKind returns all elements of the specified kind (tag).
func ElementsByKind(kind ElementKind) []Element {
	elems := doc.CallRaw("getElementsByTagName", string(kind))
	o := make([]Element, elems.Length())
	for i := range o {
		o[i] = Element{elems.Index(i)}
	}
	return o
}

// Handle registers a document event listener.
func Handle(event EventName, h HandlerFunction) {
	doc.CallRaw("addEventListener", string(event), wasm.Value(h))
}

// HandleRemove deregisters a document event listener.
func HandleRemove(event EventName, h HandlerFunction) {
	doc.CallRaw("removeEventListener", string(event), wasm.Value(h))
}

// Url returns the current navigation URL.
func Url() url.URL {
	s := location.Get("href").String()
	u, _ := url.Parse(s)
	return *u
}

func WindowHandle(event EventName, h HandlerFunction) {
	wasm.Global.CallRaw("addEventListener", string(event), wasm.Value(h))
}

func WindowHandleRemove(event EventName, h HandlerFunction) {
	wasm.Global.CallRaw("removeEventListener", string(event), wasm.Value(h))
}
