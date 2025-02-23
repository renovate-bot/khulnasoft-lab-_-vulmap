package vulmap

import (
	"context"
	"time"

	"github.com/khulnasoft-lab/gologger"
	"github.com/khulnasoft-lab/vulmap/pkg/model/types/severity"
	"github.com/khulnasoft-lab/vulmap/pkg/output"
	"github.com/khulnasoft-lab/vulmap/pkg/progress"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/hosterrorscache"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/interactsh"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/common/utils/vardump"
	"github.com/khulnasoft-lab/vulmap/pkg/protocols/headless/engine"
	"github.com/khulnasoft-lab/vulmap/pkg/templates/types"
	"github.com/khulnasoft-lab/ratelimit"
)

// TemplateSources contains template sources
// which define where to load templates from
type TemplateSources struct {
	Templates       []string // template file/directory paths
	Workflows       []string // workflow file/directory paths
	RemoteTemplates []string // remote template urls
	RemoteWorkflows []string // remote workflow urls
	TrustedDomains  []string // trusted domains for remote templates/workflows
}

// WithTemplatesOrWorkflows sets templates / workflows to use /load
func WithTemplatesOrWorkflows(sources TemplateSources) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		// by default all of these values are empty
		e.opts.Templates = sources.Templates
		e.opts.Workflows = sources.Workflows
		e.opts.TemplateURLs = sources.RemoteTemplates
		e.opts.WorkflowURLs = sources.RemoteWorkflows
		e.opts.RemoteTemplateDomainList = append(e.opts.RemoteTemplateDomainList, sources.TrustedDomains...)
		return nil
	}
}

// config contains all SDK configuration options
type TemplateFilters struct {
	Severity             string   // filter by severities (accepts CSV values of info, low, medium, high, critical)
	ExcludeSeverities    string   // filter by excluding severities (accepts CSV values of info, low, medium, high, critical)
	ProtocolTypes        string   // filter by protocol types
	ExcludeProtocolTypes string   // filter by excluding protocol types
	Authors              []string // fiter by author
	Tags                 []string // filter by tags present in template
	ExcludeTags          []string // filter by excluding tags present in template
	IncludeTags          []string // filter by including tags present in template
	IDs                  []string // filter by template IDs
	ExcludeIDs           []string // filter by excluding template IDs
	TemplateCondition    []string // DSL condition/ expression
}

// WithTemplateFilters sets template filters and only templates matching the filters will be
// loaded and executed
func WithTemplateFilters(filters TemplateFilters) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		s := severity.Severities{}
		if err := s.Set(filters.Severity); err != nil {
			return err
		}
		es := severity.Severities{}
		if err := es.Set(filters.ExcludeSeverities); err != nil {
			return err
		}
		pt := types.ProtocolTypes{}
		if err := pt.Set(filters.ProtocolTypes); err != nil {
			return err
		}
		ept := types.ProtocolTypes{}
		if err := ept.Set(filters.ExcludeProtocolTypes); err != nil {
			return err
		}
		e.opts.Authors = filters.Authors
		e.opts.Tags = filters.Tags
		e.opts.ExcludeTags = filters.ExcludeTags
		e.opts.IncludeTags = filters.IncludeTags
		e.opts.IncludeIds = filters.IDs
		e.opts.ExcludeIds = filters.ExcludeIDs
		e.opts.Severities = s
		e.opts.ExcludeSeverities = es
		e.opts.Protocols = pt
		e.opts.ExcludeProtocols = ept
		e.opts.IncludeConditions = filters.TemplateCondition
		return nil
	}
}

// InteractshOpts contains options for interactsh
type InteractshOpts interactsh.Options

// WithInteractshOptions sets interactsh options
func WithInteractshOptions(opts InteractshOpts) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("WithInteractshOptions")
		}
		optsPtr := &opts
		e.interactshOpts = (*interactsh.Options)(optsPtr)
		return nil
	}
}

// Concurrency options
type Concurrency struct {
	TemplateConcurrency         int // number of templates to run concurrently (per host in host-spray mode)
	HostConcurrency             int // number of hosts to scan concurrently  (per template in template-spray mode)
	HeadlessHostConcurrency     int // number of hosts to scan concurrently for headless templates  (per template in template-spray mode)
	HeadlessTemplateConcurrency int // number of templates to run concurrently for headless templates (per host in host-spray mode)
}

// WithConcurrency sets concurrency options
func WithConcurrency(opts Concurrency) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		e.opts.TemplateThreads = opts.TemplateConcurrency
		e.opts.BulkSize = opts.HostConcurrency
		e.opts.HeadlessBulkSize = opts.HeadlessHostConcurrency
		e.opts.HeadlessTemplateThreads = opts.HeadlessTemplateConcurrency
		return nil
	}
}

// WithGlobalRateLimit sets global rate (i.e all hosts combined) limit options
func WithGlobalRateLimit(maxTokens int, duration time.Duration) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		e.rateLimiter = ratelimit.New(context.Background(), uint(maxTokens), duration)
		return nil
	}
}

// HeadlessOpts contains options for headless templates
type HeadlessOpts struct {
	PageTimeout     int // timeout for page load
	ShowBrowser     bool
	HeadlessOptions []string
	UseChrome       bool
}

// EnableHeadless allows execution of headless templates
// *Use With Caution*: Enabling headless mode may open up attack surface due to browser usage
// and can be prone to exploitation by custom unverified templates if not properly configured
func EnableHeadlessWithOpts(hopts *HeadlessOpts) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		e.opts.Headless = true
		if hopts != nil {
			e.opts.HeadlessOptionalArguments = hopts.HeadlessOptions
			e.opts.PageTimeout = hopts.PageTimeout
			e.opts.ShowBrowser = hopts.ShowBrowser
			e.opts.UseInstalledChrome = hopts.UseChrome
		}
		if engine.MustDisableSandbox() {
			gologger.Warning().Msgf("The current platform and privileged user will run the browser without sandbox\n")
		}
		browser, err := engine.New(e.opts)
		if err != nil {
			return err
		}
		e.executerOpts.Browser = browser
		return nil
	}
}

// StatsOptions
type StatsOptions struct {
	Interval         int
	JSON             bool
	MetricServerPort int
}

// EnableStats enables Stats collection with defined interval(in sec) and callback
// Note: callback is executed in a separate goroutine
func EnableStatsWithOpts(opts StatsOptions) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("EnableStatsWithOpts")
		}
		if opts.Interval == 0 {
			opts.Interval = 5 //sec
		}
		e.opts.StatsInterval = opts.Interval
		e.enableStats = true
		e.opts.StatsJSON = opts.JSON
		e.opts.MetricsPort = opts.MetricServerPort
		return nil
	}
}

// VerbosityOptions
type VerbosityOptions struct {
	Verbose       bool // show verbose output
	Silent        bool // show only results
	Debug         bool // show debug output
	DebugRequest  bool // show request in debug output
	DebugResponse bool // show response in debug output
	ShowVarDump   bool // show variable dumps in output
}

// WithVerbosity allows setting verbosity options of (internal) vulmap engine
// and does not affect SDK output
func WithVerbosity(opts VerbosityOptions) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("WithVerbosity")
		}
		e.opts.Verbose = opts.Verbose
		e.opts.Silent = opts.Silent
		e.opts.Debug = opts.Debug
		e.opts.DebugRequests = opts.DebugRequest
		e.opts.DebugResponse = opts.DebugResponse
		if opts.ShowVarDump {
			vardump.EnableVarDump = true
		}
		return nil
	}
}

// NetworkConfig contains network config options
// ex: retries , httpx probe , timeout etc
type NetworkConfig struct {
	Timeout           int      // Timeout in seconds
	Retries           int      // Number of retries
	LeaveDefaultPorts bool     // Leave default ports for http/https
	MaxHostError      int      // Maximum number of host errors to allow before skipping that host
	TrackError        []string // Adds given errors to max host error watchlist
	DisableMaxHostErr bool     // Disable max host error optimization (Hosts are not skipped even if they are not responding)
}

// WithNetworkConfig allows setting network config options
func WithNetworkConfig(opts NetworkConfig) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("WithNetworkConfig")
		}
		e.opts.Timeout = opts.Timeout
		e.opts.Retries = opts.Retries
		e.opts.LeaveDefaultPorts = opts.LeaveDefaultPorts
		e.hostErrCache = hosterrorscache.New(opts.MaxHostError, hosterrorscache.DefaultMaxHostsCount, opts.TrackError)
		return nil
	}
}

// WithProxy allows setting proxy options
func WithProxy(proxy []string, proxyInternalRequests bool) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("WithProxy")
		}
		e.opts.Proxy = proxy
		e.opts.ProxyInternal = proxyInternalRequests
		return nil
	}
}

// WithScanStrategy allows setting scan strategy options
func WithScanStrategy(strategy string) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		e.opts.ScanStrategy = strategy
		return nil
	}
}

// OutputWriter
type OutputWriter output.Writer

// UseWriter allows setting custom output writer
// by default a mock writer is used with user defined callback
// if outputWriter is used callback will be ignored
func UseOutputWriter(writer OutputWriter) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("UseOutputWriter")
		}
		e.customWriter = writer
		return nil
	}
}

// StatsWriter
type StatsWriter progress.Progress

// UseStatsWriter allows setting a custom stats writer
// which can be used to write stats somewhere (ex: send to webserver etc)
func UseStatsWriter(writer StatsWriter) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("UseStatsWriter")
		}
		e.customProgress = writer
		return nil
	}
}

// WithTemplateUpdateCallback allows setting a callback which will be called
// when vulmap templates are outdated
// Note: Vulmap-templates are crucial part of vulmap and using outdated templates or vulmap sdk is not recommended
// as it may cause unexpected results due to compatibility issues
func WithTemplateUpdateCallback(disableTemplatesAutoUpgrade bool, callback func(newVersion string)) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("WithTemplateUpdateCallback")
		}
		e.disableTemplatesAutoUpgrade = disableTemplatesAutoUpgrade
		e.onUpdateAvailableCallback = callback
		return nil
	}
}

// WithSandboxOptions allows setting supported sandbox options
func WithSandboxOptions(allowLocalFileAccess bool, restrictLocalNetworkAccess bool) VulmapSDKOptions {
	return func(e *VulmapEngine) error {
		if e.mode == threadSafe {
			return ErrOptionsNotSupported.Msgf("WithSandboxOptions")
		}
		e.opts.AllowLocalFileAccess = allowLocalFileAccess
		e.opts.RestrictLocalNetworkAccess = restrictLocalNetworkAccess
		return nil
	}
}
