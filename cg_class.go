package goapp_generator

import (
	"fmt"
	R "goapp_commons/reflext"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/exp/slog"
)

var type2class map[reflect.Type]*Class = make(map[reflect.Type]*Class)

var name2class map[string]*Class = make(map[string]*Class)

var classes []*Class

func ForClassType(typ reflect.Type) *Class {
	return type2class[typ]
}

func ForClassName(fullname string) *Class {
	return name2class[fullname]
}

func ForClassAny(value any) *Class {
	t := reflect.TypeOf(value)
	if t == nil {
		return nil
	}
	return ForClassType(t)
}

func RegisterClass(cls *Class) {
	if cls == nil {
		panic("RegisterClass cls is nil")
	}
	if cls.protoType != nil {
		if ForClassType(cls.protoType) != nil {
			panic(fmt.Sprintf("RegisterClass type conflict '%v'", cls.protoType))
		}
		type2class[cls.protoType] = cls
	}
	if cls.fullname != "" {
		if ForClassName(cls.fullname) != nil {
			panic(fmt.Sprintf("RegisterClass name conflict '%v'", cls.fullname))
		}
		name2class[cls.fullname] = cls
	}
	classes = append(classes, cls)
}

func SelectAllClass() []*Class {
	return classes
}

type Kind int

const (
	KindValue Kind = iota
	KindObject
	KindFunc
	KindProperty
	KindMethod
	KindInterface
)

func (k Kind) String() string {
	switch k {
	case KindValue:
		return "Value"
	case KindObject:
		return "Object"
	case KindFunc:
		return "Func"
	case KindProperty:
		return "Property"
	case KindMethod:
		return "Method"
	case KindInterface:
		return "Interface"
	default:
		return "Unknow"
	}
}

type Class struct {
	kind         Kind
	proto        any
	protoType    reflect.Type
	fullname     string
	pkg          string
	name         string
	goGenPkgpath string
	goGenPkgname string
	goGenFile    string
	detail       any
}

func newProtoClass(k Kind, proto any) *Class {
	if proto == nil {
		panic("newProtoClass proto is nil")
	}
	o := &Class{
		kind: k,
	}
	o.proto = proto
	typ := reflect.TypeOf(proto)
	o.initClass(typ)

	return o
}

func newTemplateClass(k Kind, template any, fieldName string) *Class {
	tt := reflect.TypeOf(template)
	field, ok := tt.FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("%T miss Field %s", template, fieldName))
	}
	return newTypeClass(k, field.Type)
}

func newTypeClass(k Kind, typ reflect.Type) *Class {
	if typ == nil {
		panic("newTypeClass typ is nil")
	}
	o := &Class{
		kind: k,
	}
	o.proto = nil
	o.initClass(typ)

	return o
}

func (th *Class) initClass(typ reflect.Type) {
	th.protoType = typ
	th.fullname = R.TypeFullname(th.protoType)
	th.pkg = th.protoType.PkgPath()
	th.name = th.protoType.Name()

	th.WithGoGenFile(GoGenFile)
	th.WithGoGenPkg(GoGenPackPath, GoGenPackName)
}

func (th *Class) Kind() Kind {
	return th.kind
}

func (th *Class) ProtoType() reflect.Type {
	return th.protoType
}

func (th *Class) NameWithPkg() R.NameWithPkg {
	return R.NP(th.pkg, th.name)
}

func (th *Class) WithGoGenPkg(p string, n string) *Class {
	th.goGenPkgpath = p
	if n == "" {
		idx := strings.LastIndex(p, "/")
		if idx >= 0 {
			n = p[idx+1:]
		}
	}
	th.goGenPkgname = n
	return th
}

func (th *Class) WithName(n string) *Class {
	th.name = n
	return th
}

func (th *Class) WithGoGenFile(n string) *Class {
	th.goGenFile = n
	return th
}

func (th *Class) GetGoGenFilePath(dir string) (fn string, share bool) {
	var n string
	s := false
	if th.goGenFile == "" {
		n = fmt.Sprintf("%s%s", th.name, GO_GENFILE_EXT)
	} else {
		n = th.goGenFile
		if !strings.HasSuffix(n, GO_GENFILE_EXT) {
			n += GO_GENFILE_EXT
		}
		s = true
	}
	return filepath.Join(dir, n), s
}

func (th *Class) ValueDetail() *Value {
	if v, ok := th.detail.(*Value); ok {
		return v
	}
	panic(fmt.Sprintf("'%v' not CG_Value", th.detail))
}

func (th *Class) GenCode(ws *Workspace) {
	if ws.Config.GenGO {
		f := ws.BeginGoFile(th)
		defer ws.EndGoFile()

		f.Comment(fmt.Sprintf("Class: %s, Kind: %s", th.fullname, th.kind))
		if g, ok := th.detail.(GoGenerator); ok {
			g.GenGoCode(ws)
		} else {
			slog.Warn("Nothing to gen")
		}
	}
}
