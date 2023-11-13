package reporting

import (
	"github.com/khulnasoft-lab/vulmap/pkg/output"
)

// Client is a client for vulmap issue tracking module
type Client interface {
	RegisterTracker(tracker Tracker)
	RegisterExporter(exporter Exporter)
	Close()
	Clear()
	CreateIssue(event *output.ResultEvent) error
	GetReportingOptions() *Options
}
