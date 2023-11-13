package ssh

import (
	lib_ssh "github.com/khulnasoft-lab/vulmap/pkg/js/libs/ssh"

	"github.com/dop251/goja"
	"github.com/khulnasoft-lab/vulmap/pkg/js/gojs"
)

var (
	module = gojs.NewGojaModule("vulmap/ssh")
)

func init() {
	module.Set(
		gojs.Objects{
			// Functions

			// Var and consts

			// Types (value type)
			"SSHClient": func() lib_ssh.SSHClient { return lib_ssh.SSHClient{} },

			// Types (pointer type)
			"NewSSHClient": func() *lib_ssh.SSHClient { return &lib_ssh.SSHClient{} },
		},
	).Register()
}

func Enable(runtime *goja.Runtime) {
	module.Enable(runtime)
}
