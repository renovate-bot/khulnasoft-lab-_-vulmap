package protocolstate

import (
	"strings"

	"github.com/khulnasoft-lab/vulmap/pkg/catalog/config"
	errorutil "github.com/khulnasoft-lab/utils/errors"
	fileutil "github.com/khulnasoft-lab/utils/file"
)

var (
	// lfaAllowed means local file access is allowed
	lfaAllowed bool
)

// Normalizepath normalizes path and returns absolute path
// it returns error if path is not allowed
// this respects the sandbox rules and only loads files from
// allowed directories
func NormalizePath(filePath string) (string, error) {
	if lfaAllowed {
		return filePath, nil
	}
	cleaned, err := fileutil.ResolveNClean(filePath, config.DefaultConfig.GetTemplateDir())
	if err != nil {
		return "", errorutil.NewWithErr(err).Msgf("could not resolve and clean path %v", filePath)
	}
	// only allow files inside vulmap-templates directory
	// even current working directory is not allowed
	if strings.HasPrefix(cleaned, config.DefaultConfig.GetTemplateDir()) {
		return cleaned, nil
	}
	return "", errorutil.New("path %v is outside vulmap-template directory and -lfa is not enabled", filePath)
}
