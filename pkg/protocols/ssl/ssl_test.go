package ssl

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-lab/vulmap/pkg/model"
	"github.com/khulnasoft-lab/vulmap/pkg/model/types/severity"
	"github.com/khulnasoft-lab/vulmap/pkg/output"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/contextargs"
	"github.com/khulnasoft-lab/vulmap/pkg/testutils"
)

func TestSSLProtocol(t *testing.T) {
	options := testutils.DefaultOptions

	testutils.Init(options)
	templateID := "testing-ssl"
	request := &Request{
		Address: "{{Hostname}}",
	}
	executerOpts := testutils.NewMockExecuterOptions(options, &testutils.TemplateInfo{
		ID:   templateID,
		Info: model.Info{SeverityHolder: severity.Holder{Severity: severity.Low}, Name: "test"},
	})
	err := request.Compile(executerOpts)
	require.Nil(t, err, "could not compile ssl request")

	var gotEvent output.InternalEvent
	ctxArgs := contextargs.NewWithInput("scanme.sh:443")
	err = request.ExecuteWithResults(ctxArgs, nil, nil, func(event *output.InternalWrappedEvent) {
		gotEvent = event.InternalEvent
	})
	require.Nil(t, err, "could not run ssl request")
	require.NotEmpty(t, gotEvent, "could not get event items")
}

func TestGetAddress(t *testing.T) {
	address, _ := getAddress("https://scanme.sh")
	require.Equal(t, "scanme.sh:443", address, "could not get correct address")
}
