package tinyzap

import (
	"net/url"
	"path/filepath"
	"sync"

	"github.com/nc30/tinyfile"
	"go.uber.org/zap"
)

var once sync.Once

func init() {
	once.Do(func() {
		zap.RegisterSink("tinyfile", TinyFactory)
	})
}

func TinyFactory(path *url.URL) (zap.Sink, error) {
	return tinyfile.NewWriter(filepath.Join(path.Host, path.Path))
}
