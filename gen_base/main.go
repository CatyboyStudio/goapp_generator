package main

import (
	CG "goapp_generator"
	"goapp_generator/golang_classes"
)

func main() {
	golang_classes.ImportAll()

	cfg := CG.NewConfig()
	cfg.RootDirGO = "../../goapp_gpython"
	cfg.RootdirPY = "../../../pythonlib"
	cfg.GenFile = true

	ws := CG.NewWorkspace(cfg)
	CG.Run(ws)
}
