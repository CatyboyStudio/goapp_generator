package goapp_generator

import (
	R "goapp_commons/reflext"

	"golang.org/x/exp/maps"
)

var name2func map[string]*Func = make(map[string]*Func)

func RegisterFunc(f *Func) {
	if f == nil {
		panic("RegisterFunc func is nil")
	}
	fn := f.fullname.String()
	if fn != "" {
		panic("Registerfunc fullname is nil")
	}

	name2func[fn] = f
}

func ForFuncName(fullname string) *Func {
	return name2func[fullname]
}

func SelectAllFunc() []*Func {
	return maps.Values(name2func)
}

type Func struct {
	fullname R.NameWithPkg
}

func (f *Func) GenCode(ws *Workspace) {

}
