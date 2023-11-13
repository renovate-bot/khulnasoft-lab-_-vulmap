package smtp

import (
	lib_smtp "github.com/khulnasoft-lab/vulmap/pkg/js/libs/smtp"

	"github.com/dop251/goja"
	"github.com/khulnasoft-lab/vulmap/pkg/js/gojs"
)

var (
	module = gojs.NewGojaModule("vulmap/smtp")
)

func init() {
	module.Set(
		gojs.Objects{
			// Functions

			// Var and consts

			// Types (value type)
			"IsSMTPResponse": func() lib_smtp.IsSMTPResponse { return lib_smtp.IsSMTPResponse{} },
			"SMTPClient":     func() lib_smtp.SMTPClient { return lib_smtp.SMTPClient{} },

			// Types (pointer type)
			"NewIsSMTPResponse": func() *lib_smtp.IsSMTPResponse { return &lib_smtp.IsSMTPResponse{} },
			"NewSMTPClient":     func() *lib_smtp.SMTPClient { return &lib_smtp.SMTPClient{} },
		},
	).Register()
}

func Enable(runtime *goja.Runtime) {
	module.Enable(runtime)
}
