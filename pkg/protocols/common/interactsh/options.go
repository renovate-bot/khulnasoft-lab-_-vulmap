package interactsh

import (
	"time"

	"github.com/projectdiscovery/interactsh/pkg/client"
	"github.com/khulnasoft-lab/vulmap/pkg/output"
	"github.com/khulnasoft-lab/vulmap/pkg/progress"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting"
	"github.com/khulnasoft-lab/retryablehttp-go"
)

// Options contains configuration options for interactsh vulmap integration.
type Options struct {
	// ServerURL is the URL of the interactsh server.
	ServerURL string
	// Authorization is the Authorization header value
	Authorization string
	// CacheSize is the numbers of requests to keep track of at a time.
	// Older items are discarded in LRU manner in favor of new requests.
	CacheSize int
	// Eviction is the period of time after which to automatically discard
	// interaction requests.
	Eviction time.Duration
	// CooldownPeriod is additional time to wait for interactions after closing
	// of the poller.
	CooldownPeriod time.Duration
	// PollDuration is the time to wait before each poll to the server for interactions.
	PollDuration time.Duration
	// Output is the output writer for vulmap
	Output output.Writer
	// IssuesClient is a client for issue exporting
	IssuesClient reporting.Client
	// Progress is the vulmap progress bar implementation.
	Progress progress.Progress
	// Debug specifies whether debugging output should be shown for interactsh-client
	Debug bool
	// DebugRequest outputs interaction request
	DebugRequest bool
	// DebugResponse outputs interaction response
	DebugResponse bool
	// DisableHttpFallback controls http retry in case of https failure for server url
	DisableHttpFallback bool
	// NoInteractsh disables the engine
	NoInteractsh bool
	// NoColor disables printing colors for matches
	NoColor bool

	StopAtFirstMatch bool
	HTTPClient       *retryablehttp.Client
}

// DefaultOptions returns the default options for interactsh client
func DefaultOptions(output output.Writer, reporting reporting.Client, progress progress.Progress) *Options {
	return &Options{
		ServerURL:           client.DefaultOptions.ServerURL,
		CacheSize:           5000,
		Eviction:            60 * time.Second,
		CooldownPeriod:      5 * time.Second,
		PollDuration:        5 * time.Second,
		Output:              output,
		IssuesClient:        reporting,
		Progress:            progress,
		DisableHttpFallback: true,
		NoColor:             false,
	}
}
