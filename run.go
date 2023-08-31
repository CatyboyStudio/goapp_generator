package goapp_generator

import (
	"golang.org/x/exp/slog"
)

type Generator interface {
	GenGoCode(ws *Workspace)

	GenPyCode(ws *Workspace)
}

func Run(ws *Workspace) {
	cfg := ws.Config
	slog.Info("----- RUN -----")
	for _, f := range SelectAllFunc() {
		if !cfg.GenAll {
			if cfg.GenMatcher == nil || !cfg.GenMatcher(nil, f) {
				continue
			}
		}
		slog.Info("GEN", "Func", f.fullname)
		f.GenCode(ws)
	}

	for _, cls := range SelectAllClass() {
		if !cfg.GenAll {
			if cfg.GenMatcher == nil || !cfg.GenMatcher(cls, nil) {
				continue
			}
		}
		slog.Info("GEN", "Class", cls.fullname)
		cls.GenCode(ws)
	}

	ws.EndAll()
	slog.Info("----- END -----")
}
