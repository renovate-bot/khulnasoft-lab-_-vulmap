package multiproto_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/khulnasoft-lab/vulmap/pkg/catalog/config"
	"github.com/khulnasoft-lab/vulmap/pkg/catalog/disk"
	"github.com/khulnasoft-lab/vulmap/pkg/parsers"
	"github.com/khulnasoft-lab/vulmap/pkg/progress"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/contextargs"
	"github.com/khulnasoft-lab/vulmap/pkg/templates"
	"github.com/khulnasoft-lab/vulmap/pkg/testutils"
	"github.com/khulnasoft-lab/ratelimit"
	"github.com/stretchr/testify/require"
)

var executerOpts protocols.ExecutorOptions

func setup() {
	options := testutils.DefaultOptions
	testutils.Init(options)
	progressImpl, _ := progress.NewStatsTicker(0, false, false, false, 0)

	executerOpts = protocols.ExecutorOptions{
		Output:       testutils.NewMockOutputWriter(),
		Options:      options,
		Progress:     progressImpl,
		ProjectFile:  nil,
		IssuesClient: nil,
		Browser:      nil,
		Catalog:      disk.NewCatalog(config.DefaultConfig.TemplatesDirectory),
		RateLimiter:  ratelimit.New(context.Background(), uint(options.RateLimit), time.Second),
	}
	workflowLoader, err := parsers.NewLoader(&executerOpts)
	if err != nil {
		log.Fatalf("Could not create workflow loader: %s\n", err)
	}
	executerOpts.WorkflowLoader = workflowLoader
}

func TestMultiProtoWithDynamicExtractor(t *testing.T) {
	setup()
	Template, err := templates.Parse("testcases/multiprotodynamic.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.Equal(t, 2, len(Template.RequestsQueue))

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	gotresults, err := Template.Executer.Execute(contextargs.NewWithInput("blog.khulnasoft-lab.io"))
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)
}

func TestMultiProtoWithProtoPrefix(t *testing.T) {
	setup()
	Template, err := templates.Parse("testcases/multiprotowithprefix.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.Equal(t, 3, len(Template.RequestsQueue))

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	gotresults, err := Template.Executer.Execute(contextargs.NewWithInput("blog.khulnasoft-lab.io"))
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)
}
