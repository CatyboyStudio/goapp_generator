package goapp_generator

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/dave/jennifer/jen"
	"golang.org/x/exp/slog"
)

type Workspace struct {
	Config *Config

	HasUnknow bool

	class   *Class
	GoFile  *jen.File
	goFiles map[string]*jen.File
}

func NewWorkspace(c *Config) *Workspace {
	return &Workspace{
		Config:  c,
		goFiles: make(map[string]*jen.File),
	}
}

func (ws *Workspace) BeginGoFile(cls *Class) *jen.File {
	fn, share := cls.GetGoGenFilePath(ws.Config.RootDirGO)

	var f *jen.File
	if share {
		if sf, ok := ws.goFiles[fn]; ok {
			f = sf
		} else {
			f = jen.NewFilePathName(cls.goGenPkgpath, cls.goGenPkgname)
			ws.goFiles[fn] = f
		}
	} else {
		f = jen.NewFilePathName(cls.goGenPkgpath, cls.goGenPkgname)
	}
	ws.GoFile = f
	ws.class = cls

	f.ImportName(GPYTHON_PACK, "")

	return f
}

func (ws *Workspace) EndGoFile() {
	defer func() {
		ws.GoFile = nil
		ws.class = nil
		ws.HasUnknow = false
	}()

	if ws.GoFile == nil {
		return
	}

	for _, f := range ws.goFiles {
		if f == ws.GoFile {
			// shared, not save
			return
		}
	}

	var fn string
	if ws.class != nil {
		fn, _ = ws.class.GetGoGenFilePath(ws.Config.RootDirGO)
	}

	if fn != "" {
		ws.saveGoFile(fn, ws.GoFile)
	}
}

func (ws *Workspace) saveGoFile(fn string, f *jen.File) {
	slog.Info("Render", "file", fn)
	buf := bytes.NewBuffer(nil)
	var err error
	if err = f.Render(buf); err != nil {
		slog.Error("RenderFail", "error", err)
		time.Sleep(time.Millisecond)
		s := buf.String()
		fmt.Println(s)
		return
	}

	if ws.HasUnknow {
		slog.Warn("HasUnknow", "file", fn)
		return
	}
	co := buf.String()
	co = "// Code generated by goapp_generator. DO NOT EDIT.\n\n" + co

	if ws.Config.GenFile {
		slog.Info("Save", "file", fn)
		err = os.WriteFile(fn, []byte(co), 0644)
		if err != nil {
			slog.Error("SaveFail", "error", err)
		}
	} else {
		time.Sleep(time.Millisecond)
		fmt.Println("######################")
		fmt.Printf("%s\n", co)
		fmt.Println("######################")
	}
}

func (ws *Workspace) EndAll() {
	for fn, f := range ws.goFiles {
		ws.saveGoFile(fn, f)
	}
	clear(ws.goFiles)
}
