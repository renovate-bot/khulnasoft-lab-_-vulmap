package flow_test

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

func TestFlowTemplateWithIndex(t *testing.T) {
	// test
	setup()
	Template, err := templates.Parse("testcases/vulmap-flow-dns.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	input := contextargs.NewWithInput("hackerone.com")
	gotresults, err := Template.Executer.Execute(input)
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)
}

func TestFlowTemplateWithID(t *testing.T) {
	setup()
	// apart from parse->compile->execution this testcase checks support for use custom id for protocol request and invocation of
	// the same in js
	Template, err := templates.Parse("testcases/vulmap-flow-dns-id.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	target := contextargs.NewWithInput("hackerone.com")
	gotresults, err := Template.Executer.Execute(target)
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)
}

func TestFlowWithProtoPrefix(t *testing.T) {
	// test
	setup()

	// apart from parse->compile->execution this testcase checks
	// mix of custom protocol request id and index is supported in js
	// and also validates availability of protocol response variables in template context
	Template, err := templates.Parse("testcases/vulmap-flow-dns-prefix.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	input := contextargs.NewWithInput("hackerone.com")
	gotresults, err := Template.Executer.Execute(input)
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)
}

func TestFlowWithConditionNegative(t *testing.T) {
	setup()

	// apart from parse->compile->execution this testcase checks
	// if bitwise operator (&&) are properly executed and working
	Template, err := templates.Parse("testcases/condition-flow.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	input := contextargs.NewWithInput("scanme.sh")
	// expect no results and verify thant dns request is executed and http is not
	gotresults, err := Template.Executer.Execute(input)
	require.Nil(t, err, "could not execute template")
	require.False(t, gotresults)
}

func TestFlowWithConditionPositive(t *testing.T) {
	setup()

	// apart from parse->compile->execution this testcase checks
	// if bitwise operator (&&) are properly executed and working
	Template, err := templates.Parse("testcases/condition-flow.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	input := contextargs.NewWithInput("blog.khulnasoft-lab.io")
	// positive match . expect results also verify that both dns() and http() were executed
	gotresults, err := Template.Executer.Execute(input)
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)
}

func TestFlowWithNoMatchers(t *testing.T) {
	// when using conditional flow with no matchers at all
	// we implicitly assume that request was successful and internally changed the result to true (for scope of condition only)

	// testcase-1 : no matchers but contains extractor
	Template, err := templates.Parse("testcases/condition-flow-extractors.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	// positive match . expect results also verify that both dns() and http() were executed
	gotresults, err := Template.Executer.Execute(contextargs.NewWithInput("blog.khulnasoft-lab.io"))
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)

	// testcase-2 : no matchers and no extractors
	Template, err = templates.Parse("testcases/condition-flow-no-operators.yaml", nil, executerOpts)
	require.Nil(t, err, "could not parse template")

	require.True(t, Template.Flow != "", "not a flow template") // this is classifer if template is flow or not

	err = Template.Executer.Compile()
	require.Nil(t, err, "could not compile template")

	// positive match . expect results also verify that both dns() and http() were executed
	gotresults, err = Template.Executer.Execute(contextargs.NewWithInput("blog.khulnasoft-lab.io"))
	require.Nil(t, err, "could not execute template")
	require.True(t, gotresults)

}
