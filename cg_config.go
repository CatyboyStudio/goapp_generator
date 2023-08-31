package goapp_generator

type GenMatcher func(*Class, *Func) bool

type Config struct {
	RootDirGO string
	RootdirPY string

	GenFile bool
	GenGO   bool
	GenPY   bool

	GenMatcher GenMatcher
}

func NewConfig() *Config {
	return &Config{
		GenFile: true,
		GenGO:   true,
		GenPY:   true,
	}
}

func (c *Config) FilterByGenPkg(p, n, fn string) {
	c.GenMatcher = func(c *Class, f *Func) bool {
		if c != nil {
			if p != "*" {
				if c.goGenPkgpath != p {
					return false
				}
			}
			if n != "*" {
				if c.goGenPkgname != n {
					return false
				}
			}
			if fn != "*" {
				if c.goGenFile != fn {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}
}
