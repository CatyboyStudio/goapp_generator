package goapp_generator

import . "goapp_commons"

type CG_Value struct {
	class *CG_Class

	// Python
	p2g *NameWithPkg // func(v py.Object) type
	g2p *NameWithPkg // func(v type) py.Object

	expandAsList bool
	expandAsMap  bool

	// String Convert
	s2g *NameWithPkg // func(v string) type
	g2s *NameWithPkg // func(v type) string | default G2S_Any(v ang) string

	// InspectorEditor
	inspectorEditor string
}

func NewValue(proto any) (*CG_Value, *CG_Class) {
	cls := newProtoClass(KindValue, proto)
	self := &CG_Value{
		class: cls,
	}
	cls.detail = self
	return self, cls
}

func (this *CG_Value) Class() *CG_Class {
	return this.class
}
