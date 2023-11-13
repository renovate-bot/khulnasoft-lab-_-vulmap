package goconsole

import (
	lib_goconsole "github.com/khulnasoft-lab/vulmap/pkg/js/libs/goconsole"

	"github.com/dop251/goja"
	"github.com/khulnasoft-lab/vulmap/pkg/js/gojs"
)

var (
	module = gojs.NewGojaModule("vulmap/goconsole")
)

func init() {
	module.Set(
		gojs.Objects{
			// Functions
			"NewGoConsolePrinter": lib_goconsole.NewGoConsolePrinter,

			// Var and consts

			// Types (value type)
			"GoConsolePrinter": func() lib_goconsole.GoConsolePrinter { return lib_goconsole.GoConsolePrinter{} },

			// Types (pointer type)
		},
	).Register()
}

func Enable(runtime *goja.Runtime) {
	module.Enable(runtime)
}
