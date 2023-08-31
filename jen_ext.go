package goapp_generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

func NewSourceCode(ws *Workspace, sc string) *jen.Statement {
	id := ws.GUID()
	ph := fmt.Sprintf("// %s", id)
	ws.AddHandler(func(s string) string {
		return strings.ReplaceAll(s, ph, sc)
	})
	return jen.Comment(ph)
}

func ReturnUseConvert(pkg, typ, varn string) *jen.Statement {
	return jen.ReturnFunc(func(g *jen.Group) {
		g.Qual(pkg, typ).Parens(jen.Id(varn))
	})
}

func AssertType(v1, ok, v2 string, typ jen.Code) *jen.Statement {
	// v1, ok := v2.(typ)
	// if ok {}
	return jen.List(jen.Id(v1), jen.Id(ok)).Op(":=").Id(v2).Assert(typ)
}
