package goapp_generator

type Config struct {
	RootDirGO string
	RootdirPY string

	GenFile bool
	GenGO   bool
	GenPY   bool
	GenAll  bool

	GenMatcher func(*Class, *Func) bool
}

func NewConfig() *Config {
	return &Config{
		GenFile: true,
		GenGO:   true,
		GenPY:   true,
		GenAll:  true,
	}
}
