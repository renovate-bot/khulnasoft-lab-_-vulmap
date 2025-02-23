package reporting

import (
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/exporters/es"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/exporters/jsonexporter"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/exporters/jsonl"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/exporters/markdown"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/exporters/sarif"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/exporters/splunk"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/trackers/github"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/trackers/gitlab"
	"github.com/khulnasoft-lab/vulmap/pkg/reporting/trackers/jira"
	"github.com/khulnasoft-lab/retryablehttp-go"
)

// Options is a configuration file for vulmap reporting module
type Options struct {
	// AllowList contains a list of allowed events for reporting module
	AllowList *Filter `yaml:"allow-list"`
	// DenyList contains a list of denied events for reporting module
	DenyList *Filter `yaml:"deny-list"`
	// GitHub contains configuration options for GitHub Issue Tracker
	GitHub *github.Options `yaml:"github"`
	// GitLab contains configuration options for GitLab Issue Tracker
	GitLab *gitlab.Options `yaml:"gitlab"`
	// Jira contains configuration options for Jira Issue Tracker
	Jira *jira.Options `yaml:"jira"`
	// MarkdownExporter contains configuration options for Markdown Exporter Module
	MarkdownExporter *markdown.Options `yaml:"markdown"`
	// SarifExporter contains configuration options for Sarif Exporter Module
	SarifExporter *sarif.Options `yaml:"sarif"`
	// ElasticsearchExporter contains configuration options for Elasticsearch Exporter Module
	ElasticsearchExporter *es.Options `yaml:"elasticsearch"`
	// SplunkExporter contains configuration options for splunkhec Exporter Module
	SplunkExporter *splunk.Options `yaml:"splunkhec"`
	// JSONExporter contains configuration options for JSON Exporter Module
	JSONExporter *jsonexporter.Options `yaml:"json"`
	// JSONLExporter contains configuration options for JSONL Exporter Module
	JSONLExporter *jsonl.Options `yaml:"jsonl"`

	HttpClient *retryablehttp.Client `yaml:"-"`
}
