package golang_classes

import (
	CG "goapp_generator"

	"github.com/dave/jennifer/jen"
)

type types struct {
	ErrorType error
}

func Bool() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(false)
	o.WithGenType(jen.Bool()).
		Use_G2P(CG.UseConverter(CG.GPYTHON_PACK, "Bool")).
		Use_P2G(CG.UseSourceCode(`
	val, err := py.MakeBool(v)
	if err != nil {
		PyLogWarn("P2G_Bool fail %v: %v\n", v, err)
		return false
	}
	return bool(val.(py.Bool))
	`)).
		Use_G2P(CG.UseConverter(CG.GPYTHON_PACK, "NewBool")).
		SkipGenAny().
		Finish()
	return o, cls
}

func Int() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(int(0))
	o.WithGenType(jen.Int()).
		WithPyType(jen.Qual(CG.GPYTHON_PACK, "Int")).
		Use_P2G(CG.UseSourceCode(`
	val, err := py.MakeGoInt(v)
	if err != nil {
		PyLogWarn("P2G_Int fail %v: %v\n", v, err)
		return 0
	}
	return val
	`)).
		// Use_G2P(CG.UseConverter(CG.GPYTHON_PACK, "Int")).
		SkipGenAny().
		Finish()
	return o, cls
}

func Int64() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(int64(0))
	o.WithGenType(jen.Int64()).
		WithPyType(jen.Qual(CG.GPYTHON_PACK, "Int")).
		Use_P2G(CG.UseSourceCode(`
	val, err := py.MakeGoInt64(v)
	if err != nil {
		PyLogWarn("P2G_Int64 fail %v: %v\n", v, err)
		return 0
	}
	return val
	`)).
		SkipGenAny().
		Finish()
	return o, cls
}

func Float32() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(float32(0))
	o.WithGenType(jen.Float32()).
		Use_G2P(CG.UseConverter(CG.GPYTHON_PACK, "Float")).
		Use_P2G(CG.UseSourceCode(`
	return float32(P2G_Float64(v))
	`)).
		SkipGenAny().
		Finish()
	return o, cls
}

func Float() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(float64(0))
	o.WithGenType(jen.Float64()).
		Use_G2P(CG.UseConverter(CG.GPYTHON_PACK, "Float")).
		Use_P2G(CG.UseSourceCode(`
	val, err := py.MakeFloat(v)
	if err != nil {
		PyLogWarn("P2G_Float fail %v: %v\n", v, err)
		return 0
	}
	return float64(val.(py.Float))
	`)).
		SkipGenAny().
		Finish()
	return o, cls
}

func String() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue("")
	o.WithGenType(jen.String()).
		Use_G2P(CG.UseConverter(CG.GPYTHON_PACK, "String")).
		Use_P2G(CG.UseSourceCode(`
	if v, err := py.StrAsString(v); err == nil {
		return v
	}
	return ""
	`)).
		SkipGenAny().
		Finish()
	return o, cls
}

func Error() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValueByTemplate(types{}, "ErrorType")
	o.WithGenType(jen.Error()).
		Use_P2G(func(th *CG.Value, ws *CG.Workspace) *jen.Statement {
			s := CG.NewSourceCode(ws, `
	if o, ok := v.(*py.Exception); ok {
		return o
	}
	val := P2G_String(v)
			`)
			return s.Return(jen.Qual("errors", "New").Params(jen.Id("val")))
		}).
		Use_G2P(CG.UseSourceCode(`
	if v == nil {
		return py.None
	}
	return py.ExceptionNewf(py.RuntimeError, "%v", v)
	`)).
		WithPyType(jen.Op("*").Qual(CG.GPYTHON_PACK, "Exception")).
		Finish()
	return o, cls
}
