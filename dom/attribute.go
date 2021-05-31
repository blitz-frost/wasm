package dom

type StringAttr struct {
	owner *Base
	name  string
}

func (x StringAttr) Get() string {
	return x.owner.Get(x.name).String()
}

func (x StringAttr) Set(a string) {
	x.owner.Set(x.name, a)
}
