package goapp_generator

import (
	"fmt"
	"reflect"
)

type NameWithPkg struct {
	PKG  string
	Name string
}

func NP(pkg, name string) NameWithPkg {
	return NameWithPkg{
		PKG:  pkg,
		Name: name,
	}
}

func NPT(t reflect.Type) NameWithPkg {
	return NameWithPkg{t.PkgPath(), t.Name()}
}

func (o NameWithPkg) IsNil() bool {
	return o.PKG == "" && o.Name == ""
}

func (this NameWithPkg) String() string {
	if this.PKG == "" {
		return this.Name
	}
	return fmt.Sprintf("%s.%s", this.PKG, this.Name)
}

func TypeFullname(typ reflect.Type) string {
	o := NPT(typ)
	if o.IsNil() {
		s := fmt.Sprintf("%v", typ)
		return s
	}
	return o.String()
}

var type2class map[reflect.Type]*CG_Class = make(map[reflect.Type]*CG_Class)

var name2class map[string]*CG_Class = make(map[string]*CG_Class)

type CG_Kind int

const (
	KindValue = iota
	KindObject
	KindFunc
	KindProperty
	KindMethod
	KindInterface
)

type CG_Class struct {
	kind      CG_Kind
	proto     any
	protoType reflect.Type
	fullname  string
	pkg       string
	name      string
	detail    any
}

func newProtoClass(k CG_Kind, proto any) *CG_Class {
	if proto == nil {
		panic("newProtoClass proto is nil")
	}
	o := &CG_Class{
		kind: k,
	}
	o.proto = proto
	o.protoType = reflect.TypeOf(proto)
	o.fullname = TypeFullname(o.protoType)
	o.pkg = o.protoType.PkgPath()
	o.name = o.protoType.Name()

	return o
}

func newTypeClass(k CG_Kind, typ reflect.Type) *CG_Class {
	if typ == nil {
		panic("newTypeClass typ is nil")
	}
	o := &CG_Class{
		kind: k,
	}
	o.proto = nil
	o.protoType = typ
	o.fullname = TypeFullname(o.protoType)
	o.pkg = o.protoType.PkgPath()
	o.name = o.protoType.Name()

	return o
}

func (this *CG_Class) Kind() CG_Kind {
	return this.kind
}

func (this *CG_Class) NameWithPkg() NameWithPkg {
	return NP(this.pkg, this.name)
}

func (this *CG_Class) WithPkg(n string) *CG_Class {
	this.pkg = n
	return this
}

func (this *CG_Class) WithName(n string) *CG_Class {
	this.name = n
	return this
}

func (this *CG_Class) ValueDetail() *CG_Value {
	if v, ok := this.detail.(*CG_Value); ok {
		return v
	}
	panic(fmt.Sprintf("'%v' not CG_Value", this.detail))
}

func ForClassType(typ reflect.Type) *CG_Class {
	return type2class[typ]
}

func ForClassName(fullname string) *CG_Class {
	return name2class[fullname]
}

func ForClassAny(value any) *CG_Class {
	t := reflect.TypeOf(value)
	if t == nil {
		return nil
	}
	return ForClassType(t)
}

func RegisterClass(cls *CG_Class) {
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
}
