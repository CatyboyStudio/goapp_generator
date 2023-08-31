package golang_classes

import CG "goapp_generator"

const (
	BUILDIN = "buildin"
)

func ImportAll() {

	CG.GoGenPackPath = CG.APPGPY_PACK
	CG.GoGenPackName = ""

	CG.GoGenFile = BUILDIN

	vlist := []func() (*CG.Value, *CG.Class){
		Bool, Int, Int64, Float32, Float, String, Error,
	}

	for _, f := range vlist {
		_, cls := f()
		CG.RegisterClass(cls)
	}
}
