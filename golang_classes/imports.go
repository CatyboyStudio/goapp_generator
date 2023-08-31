package golang_classes

import CG "goapp_generator"

func ImportAll() {

	CG.GoGenPackPath = "github.com/CatyboyStudio/goapp_gpython"
	CG.GoGenPackName = ""

	CG.GoGenFile = "buildin"

	vlist := []func() (*CG.Value, *CG.Class){
		Bool, Int, Int64, Float, String, Error,
	}

	for _, f := range vlist {
		_, cls := f()
		CG.RegisterClass(cls)
	}
}
