package goapp_generator

import (
	"fmt"
	R "goapp_commons/reflext"
	"strings"

	"github.com/dave/jennifer/jen"
)

var _ Generator = (*Value)(nil)

type Value struct {
	class *Class

	genName string
	genType *jen.Statement

	// Python
	p2g     *R.NameWithPkg // func(v py.Object) type
	g2p     *R.NameWithPkg // func(v type) py.Object
	buildin bool

	expandAsList bool
	expandAsMap  bool
}

func NewValue(proto any) (*Value, *Class) {
	cls := newProtoClass(KindValue, proto)
	self := &Value{
		class: cls,
	}
	cls.detail = self

	return self, cls
}

func (th *Value) WithGenName(n string) *Value {
	th.genName = n
	return th
}

func (th *Value) GenName() string {
	if th.genName == "" {
		return strings.ToUpper(th.class.name[:1]) + th.class.name[1:]
	}
	return th.genName
}

func (th *Value) WithGenType(c *jen.Statement) *Value {
	th.genType = c
	return th
}

func (th *Value) GenType() *jen.Statement {
	if th.genType != nil {
		return th.genType.Clone()
	}
	return jen.Qual(th.class.pkg, th.class.name)
}

func (th *Value) WithP2G(np R.NameWithPkg) *Value {
	v := np
	th.p2g = &v
	return th
}

func (th *Value) WithG2P(np R.NameWithPkg) *Value {
	v := np
	th.g2p = &v
	return th
}

func (th *Value) WithBuildin() *Value {
	th.buildin = true
	return th
}

func (th *Value) WithExpandAsList() *Value {
	th.expandAsList = true
	return th
}

func (th *Value) WithExpandAsMap() *Value {
	th.expandAsMap = true
	return th
}

func (th *Value) Finish() *Value {
	return th
}

func (th *Value) Class() *Class {
	return th.class
}

func (th *Value) P2G() R.NameWithPkg {
	if th.p2g != nil {
		return *th.p2g
	} else {
		n := fmt.Sprintf("P2G_%s", th.GenName())
		return R.NP(th.class.goGenPkgpath, n)
	}
}

func (th *Value) G2P() R.NameWithPkg {
	if th.p2g != nil {
		return *th.p2g
	} else {
		n := fmt.Sprintf("G2P_%s", th.GenName())
		return R.NP(th.class.goGenPkgpath, n)
	}
}

func (th *Value) GenGoCode(ws *Workspace) {
	f := ws.GoFile

	if true {
		np := th.P2G()
		if th.buildin {
			// var _ func(v py.Object) type = XXX
			f.Var().Id("_").Func().Params(
				jen.Id("v").Qual(GPYTHON_PACK, "Object"),
			).Add(th.GenType()).
				Op("=").Qual(np.PKG, np.Name)
		} else {
			// func(v py.Object) type
			f.Func().Qual(np.PKG, np.Name).Params(
				jen.Id("v").Qual(GPYTHON_PACK, "Object"),
			).Add(th.GenType()).Block(
			//
			)
		}
	}

	if true {
		np := th.G2P()
		if th.buildin {
			// var _ func(v type) py.Object = XXX
			f.Var().Id("_").Func().Params(
				jen.Id("v").Add(th.GenType()),
			).Qual(GPYTHON_PACK, "Object").
				Op("=").Qual(np.PKG, np.Name)
		} else {
			// func(v type) py.Object
			f.Func().Qual(np.PKG, np.Name).Params(
				jen.Id("v").Add(th.GenType()),
			).Qual(GPYTHON_PACK, "Object").Block(
			//
			)
		}
	}
}

func (th *Value) GenPyCode(ws *Workspace) {
}
