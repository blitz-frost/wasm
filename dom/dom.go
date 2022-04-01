package dom

import (
	"errors"
	"net/url"
	"syscall/js"
)

var (
	window   = js.Global()
	console  = window.Get("console")
	doc      = window.Get("document")
	location = window.Get("location")
)

// GetId returns the element with the given ID in the document.
// Returns an error if the ID doesn't exist.
func GetId(id string) (Element, error) {
	elem := doc.Call("getElementById", id)
	if elem.IsNull() {
		return Element{}, errors.New(id + " not found")
	}
	return Element{elem}, nil
}

// GetKind returns all elements of the specified kind (tag).
func GetKind(kind ElementKind) []Element {
	elems := doc.Call("getElementsByTagName", string(kind))
	o := make([]Element, elems.Length())
	for i := range o {
		o[i] = Element{elems.Index(i)}
	}
	return o
}

func GetUrl() url.URL {
	s := location.Get("href").String()
	u, _ := url.Parse(s)
	return *u
}

func Handle(event EventName, h Handler) {
	doc.Call("addEventListener", string(event), h.f)
}

func HandleRemove(event EventName, h Handler) {
	doc.Call("removeEventListener", string(event), h.f)
}

/*
//TODO update along with jsconv package
// Log wraps the standard package fmt.Println.
// If a is a syscall/js.Wrapper (is or can convert itself to a JS value), then it will be passed to the browser console for formatting.
func Log(a interface{}) {
	if jsw, ok := a.(js.Wrapper); ok {
		console.Call("log", jsw.JSValue())
		return
	}
	fmt.Println(a)
}
*/
