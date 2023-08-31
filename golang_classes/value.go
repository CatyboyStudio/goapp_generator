package golang_classes

import (
	"errors"
	CG "goapp_generator"

	"github.com/dave/jennifer/jen"
)

func Bool() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(false)
	o.WithGenType(jen.Bool()).
		WithBuildin().
		Finish()
	return o, cls
}

func Int() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(int(0))
	o.WithGenType(jen.Int()).
		WithBuildin().
		Finish()
	return o, cls
}

func Int64() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(int64(0))
	o.WithGenType(jen.Int64()).
		WithBuildin().
		Finish()
	return o, cls
}

func Float() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(float64(0))
	o.WithGenType(jen.Float64()).
		WithBuildin().
		Finish()
	return o, cls
}

func String() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue("")
	o.WithGenType(jen.String()).
		WithBuildin().
		Finish()
	return o, cls
}

func Error() (*CG.Value, *CG.Class) {
	o, cls := CG.NewValue(errors.New(""))
	o.WithGenType(jen.Error()).
		WithGenName("Error").
		WithBuildin().
		Finish()
	return o, cls
}
