package disk

import (
	"io"
	"io/fs"
	"os"

	"github.com/khulnasoft-lab/vulmap/pkg/catalog/config"
)

// DiskCatalog is a template catalog helper implementation based on disk
type DiskCatalog struct {
	templatesDirectory string
	templatesFS        fs.FS // TODO: Refactor to use this
}

// NewCatalog creates a new Catalog structure using provided input items
// using disk based items
func NewCatalog(directory string) *DiskCatalog {
	catalog := &DiskCatalog{templatesDirectory: directory}
	if directory != "" {
		catalog.templatesFS = os.DirFS(directory)
	} else {
		catalog.templatesFS = os.DirFS(config.DefaultConfig.GetTemplateDir())
	}
	return catalog
}

// OpenFile opens a file and returns an io.ReadCloser to the file.
// It is used to read template and payload files based on catalog responses.
func (d *DiskCatalog) OpenFile(filename string) (io.ReadCloser, error) {
	file, err := os.Open(filename)
	if err != nil {
		if file, errx := os.Open(BackwardsCompatiblePaths(d.templatesDirectory, filename)); errx == nil {
			return file, nil
		}
	}
	return file, err
}
