package goapp_generator

import (
	"fmt"
	R "goapp_commons/reflext"
	"strings"

	"github.com/dave/jennifer/jen"
)

type ValueGoGenFunc func(th *Value, ws *Workspace) *jen.Statement

var _ GoGenerator = (*Value)(nil)

type Value struct {
	class *Class

	genName string
	genType *jen.Statement
	pyType  *jen.Statement

	// GPython
	p2gGenF    ValueGoGenFunc
	g2pGenF    ValueGoGenFunc
	p2aGenF    ValueGoGenFunc
	a2pGenF    ValueGoGenFunc
	skipGenAny bool

	expandAsList bool
	expandAsMap  bool
}

func NewValue(proto any) (*Value, *Class) {
	cls := newProtoClass(KindValue, proto)
	return NewValueByClass(cls)
}

func NewValueByTemplate(template any, fieldName string) (*Value, *Class) {
	cls := newTemplateClass(KindValue, template, fieldName)
	return NewValueByClass(cls)
}

func NewValueByClass(cls *Class) (*Value, *Class) {
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

func (th *Value) WithPyType(s *jen.Statement) *Value {
	th.pyType = s
	return th
}

func (th *Value) PyType() *jen.Statement {
	return th.pyType
}

func (th *Value) Use_P2G(f ValueGoGenFunc) *Value {
	th.p2gGenF = f
	return th
}

func (th *Value) Use_G2P(f ValueGoGenFunc) *Value {
	th.g2pGenF = f
	return th
}

func (th *Value) Use_P2A(f ValueGoGenFunc) *Value {
	th.p2aGenF = f
	return th
}

func (th *Value) Use_A2P(f ValueGoGenFunc) *Value {
	th.a2pGenF = f
	return th
}

func (th *Value) SkipGenAny() *Value {
	th.skipGenAny = true
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
	n := fmt.Sprintf("P2G_%s", th.GenName())
	return R.NP(th.class.goGenPkgpath, n)
}

func (th *Value) G2P() R.NameWithPkg {
	n := fmt.Sprintf("G2P_%s", th.GenName())
	return R.NP(th.class.goGenPkgpath, n)
}

func (th *Value) P2A() R.NameWithPkg {
	n := fmt.Sprintf("P2A_%s", th.GenName())
	return R.NP(th.class.goGenPkgpath, n)
}

func (th *Value) A2P() R.NameWithPkg {
	n := fmt.Sprintf("A2P_%s", th.GenName())
	return R.NP(th.class.goGenPkgpath, n)
}

func (th *Value) GenGoCode(ws *Workspace) {
	f := ws.GoFile.jfile

	if true {
		np := th.P2G()
		// func(v py.Object) type
		f.Func().Qual(np.PKG, np.Name).Params(
			jen.Id("v").Qual(GPYTHON_PACK, "Object"),
		).Add(th.GenType()).Block(
			*th.genP2GCode(ws)...,
		)
	}

	if true {
		np := th.G2P()
		// func(v type) py.Object
		f.Func().Qual(np.PKG, np.Name).Params(
			jen.Id("v").Add(th.GenType()),
		).Qual(GPYTHON_PACK, "Object").Block(
			*th.genG2PCode(ws)...,
		)
	}

	if !th.skipGenAny && (th.pyType != nil || th.p2aGenF != nil) {
		np := th.P2A()
		// func(v py.Object) (any, bool)
		f.Func().Qual(np.PKG, np.Name).Params(
			jen.Id("v").Qual(GPYTHON_PACK, "Object"),
		).Params(jen.Any(), jen.Bool()).Block(
			*th.genP2ACode(ws)...,
		)
		ic := jen.Qual(APPGPY_PACK, "RegisterToAny").Params(jen.Qual(np.PKG, np.Name))
		ws.AddInitCode(ic)
	}

	if !th.skipGenAny {
		np := th.A2P()
		// func(v any) (py.Object, bool)
		f.Func().Qual(np.PKG, np.Name).Params(
			jen.Id("v").Any(),
		).Params(jen.Qual(GPYTHON_PACK, "Object"), jen.Bool()).Block(
			*th.genA2PCode(ws)...,
		)
		ic := jen.Qual(APPGPY_PACK, "RegisterAnyTo").Params(jen.Qual(np.PKG, np.Name))
		ws.AddInitCode(ic)
	}
}

func (th *Value) genP2GCode(ws *Workspace) *jen.Statement {
	if th.p2gGenF != nil {
		return th.p2gGenF(th, ws)
	}
	return jen.Comment("TODO")
}

func (th *Value) genG2PCode(ws *Workspace) *jen.Statement {
	if th.g2pGenF != nil {
		return th.g2pGenF(th, ws)
	}
	if th.pyType != nil {
		return jen.Return(th.pyType.Clone().Params(jen.Id("v")))
	}
	return jen.Comment("TODO")
}

func (th *Value) genP2ACode(ws *Workspace) *jen.Statement {
	if th.p2aGenF != nil {
		return th.p2aGenF(th, ws)
	}
	/* assert
	if ok {
		return P2G_xxx(v0)
	}
	*/
	pyt := th.pyType.Clone()
	s := &jen.Statement{
		AssertType("v0", "ok", "v", pyt),
		jen.If(jen.Id("ok")).BlockFunc(func(g *jen.Group) {
			np := th.P2G()
			g.Return(jen.Qual(np.PKG, np.Name).Params(jen.Id("v0")), jen.True())
		}),
		jen.Return(jen.Nil(), jen.False()),
	}
	return s
}

func (th *Value) genA2PCode(ws *Workspace) *jen.Statement {
	if th.a2pGenF != nil {
		return th.a2pGenF(th, ws)
	}
	/* assert
	if ok {
		return G2P_xxx(v0)
	}
	*/
	s := &jen.Statement{
		AssertType("v0", "ok", "v", th.GenType()),
		jen.If(jen.Id("ok")).BlockFunc(func(g *jen.Group) {
			np := th.G2P()
			g.Return(jen.Qual(np.PKG, np.Name).Params(jen.Id("v0")), jen.True())
		}),
		jen.Return(jen.Qual(GPYTHON_PACK, "None"), jen.False()),
	}
	return s
}

func UseSourceCode(s string) ValueGoGenFunc {
	return func(th *Value, ws *Workspace) *jen.Statement {
		return NewSourceCode(ws, s)
	}
}

func UseConverter(p, n string) ValueGoGenFunc {
	return func(th *Value, ws *Workspace) *jen.Statement {
		return ReturnUseConvert(p, n, "v")
	}
}
