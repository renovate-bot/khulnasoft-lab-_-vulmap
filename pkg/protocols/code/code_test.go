//go:build linux || darwin

package code

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-lab/vulmap/pkg/model"
	"github.com/khulnasoft-lab/vulmap/pkg/model/types/severity"
	"github.com/khulnasoft-lab/vulmap/pkg/output"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/contextargs"
	"github.com/khulnasoft-lab/vulmap/pkg/testutils"
)

func TestCodeProtocol(t *testing.T) {
	options := testutils.DefaultOptions

	testutils.Init(options)
	templateID := "testing-code"
	request := &Request{
		Engine: []string{"sh"},
		Source: "echo test",
	}
	executerOpts := testutils.NewMockExecuterOptions(options, &testutils.TemplateInfo{
		ID:   templateID,
		Info: model.Info{SeverityHolder: severity.Holder{Severity: severity.Low}, Name: "test"},
	})
	err := request.Compile(executerOpts)
	require.Nil(t, err, "could not compile code request")

	var gotEvent output.InternalEvent
	ctxArgs := contextargs.NewWithInput("")
	err = request.ExecuteWithResults(ctxArgs, nil, nil, func(event *output.InternalWrappedEvent) {
		gotEvent = event.InternalEvent
	})
	require.Nil(t, err, "could not run code request")
	require.NotEmpty(t, gotEvent, "could not get event items")
}
