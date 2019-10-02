package driver

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-index"
	"os"
	"path/filepath"
)

func init() {
	dr := &DirectoryDriver{}
	index.Register("directory", dr)
}

type DirectoryDriver struct {
	index.Driver
}

func (d *DirectoryDriver) Open(uri string) error {
	return nil
}

func (d *DirectoryDriver) IndexURI(ctx context.Context, index_cb index.IndexerFunc, uri string) error {

	abs_path, err := filepath.Abs(uri)

	if err != nil {
		return err
	}

	crawl_cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		if info.IsDir() {
			return nil
		}

		fh, err := readerFromPath(path)

		if err != nil {
			return err
		}

		defer fh.Close()

		ctx = index.AssignPathContext(ctx, path)
		return index_cb(ctx, fh)
	}

	c := crawl.NewCrawler(abs_path)
	return c.Crawl(crawl_cb)
}
