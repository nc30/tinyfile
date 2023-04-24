package tinyzap

import (
	"net/url"
	"path/filepath"

	"github.com/nc30/tinyfile"
	"go.uber.org/zap"
)

func TinyFactory(path *url.URL) (zap.Sink, error) {
	return tinyfile.NewWriter(filepath.Join(path.Host, path.Path))
}
