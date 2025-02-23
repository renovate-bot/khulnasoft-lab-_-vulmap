package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	"github.com/khulnasoft-lab/goflags"
	"github.com/khulnasoft-lab/vulmap/pkg/catalog/config"
	"github.com/khulnasoft-lab/vulmap/pkg/catalog/disk"
	"github.com/khulnasoft-lab/vulmap/pkg/catalog/loader"
	"github.com/khulnasoft-lab/vulmap/pkg/core"
	"github.com/khulnasoft-lab/vulmap/pkg/core/inputs"
	"github.com/khulnasoft-lab/vulmap/pkg/output"
	"github.com/khulnasoft-lab/vulmap/pkg/parsers"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/contextargs"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/hosterrorscache"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/interactsh"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/protocolinit"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/protocolstate"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting"
	"github.com/khulnasoft-lab/vulmap/pkg/testutils"
	"github.com/khulnasoft-lab/vulmap/pkg/types"
	"github.com/khulnasoft-lab/ratelimit"
)

var libraryTestcases = []TestCaseInfo{
	{Path: "library/test.yaml", TestCase: &goIntegrationTest{}},
	{Path: "library/test.json", TestCase: &goIntegrationTest{}},
}

type goIntegrationTest struct{}

// Execute executes a test case and returns an error if occurred
//
// Execute the docs at ../DESIGN.md if the code stops working for integration.
func (h *goIntegrationTest) Execute(templatePath string) error {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "This is test matcher text")
		if strings.EqualFold(r.Header.Get("test"), "vulmap") {
			fmt.Fprintf(w, "This is test headers matcher text")
		}
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	results, err := executeVulmapAsLibrary(templatePath, ts.URL)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}

// executeVulmapAsLibrary contains an example
func executeVulmapAsLibrary(templatePath, templateURL string) ([]string, error) {
	cache := hosterrorscache.New(30, hosterrorscache.DefaultMaxHostsCount, nil)
	defer cache.Close()

	mockProgress := &testutils.MockProgressClient{}
	reportingClient, err := reporting.New(&reporting.Options{}, "")
	if err != nil {
		return nil, err
	}
	defer reportingClient.Close()

	outputWriter := testutils.NewMockOutputWriter()
	var results []string
	outputWriter.WriteCallback = func(event *output.ResultEvent) {
		results = append(results, fmt.Sprintf("%v\n", event))
	}

	defaultOpts := types.DefaultOptions()
	_ = protocolstate.Init(defaultOpts)
	_ = protocolinit.Init(defaultOpts)

	defaultOpts.Templates = goflags.StringSlice{templatePath}
	defaultOpts.ExcludeTags = config.ReadIgnoreFile().Tags

	interactOpts := interactsh.DefaultOptions(outputWriter, reportingClient, mockProgress)
	interactClient, err := interactsh.New(interactOpts)
	if err != nil {
		return nil, errors.Wrap(err, "could not create interact client")
	}
	defer interactClient.Close()

	home, _ := os.UserHomeDir()
	catalog := disk.NewCatalog(path.Join(home, "vulmap-templates"))
	ratelimiter := ratelimit.New(context.Background(), 150, time.Second)
	defer ratelimiter.Stop()
	executerOpts := protocols.ExecutorOptions{
		Output:          outputWriter,
		Options:         defaultOpts,
		Progress:        mockProgress,
		Catalog:         catalog,
		IssuesClient:    reportingClient,
		RateLimiter:     ratelimiter,
		Interactsh:      interactClient,
		HostErrorsCache: cache,
		Colorizer:       aurora.NewAurora(true),
		ResumeCfg:       types.NewResumeCfg(),
	}
	engine := core.New(defaultOpts)
	engine.SetExecuterOptions(executerOpts)

	workflowLoader, err := parsers.NewLoader(&executerOpts)
	if err != nil {
		log.Fatalf("Could not create workflow loader: %s\n", err)
	}
	executerOpts.WorkflowLoader = workflowLoader

	store, err := loader.New(loader.NewConfig(defaultOpts, catalog, executerOpts))
	if err != nil {
		return nil, errors.Wrap(err, "could not create loader")
	}
	store.Load()

	input := &inputs.SimpleInputProvider{Inputs: []*contextargs.MetaInput{{Input: templateURL}}}
	_ = engine.Execute(store.Templates(), input)
	engine.WorkPool().Wait() // Wait for the scan to finish

	return results, nil
}
