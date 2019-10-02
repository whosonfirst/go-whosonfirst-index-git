package git

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index"	
	"github.com/whosonfirst/go-whosonfirst-index/driver"
	gogit "gopkg.in/src-d/go-git.v4"
	_ "gopkg.in/src-d/go-git.v4/plumbing/protocol/packp/sideband"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {

	rd := &driver.RepoDriver{}

	dr := &GitDriver{
		repo_driver: rd,
	}

	index.Register("git", dr)
}

type GitDriver struct {
	index.Driver
	repo_driver index.Driver
}

func (d *GitDriver) Open(uri string) error {
	return d.repo_driver.Open(uri)
}

func (d *GitDriver) IndexURI(ctx context.Context, index_cb index.IndexerFunc, uri string) error {

	repo_name := filepath.Base(uri)

	tempdir, err := ioutil.TempDir("", repo_name)

	if err != nil {
		return err
	}

	defer os.RemoveAll(tempdir)

	/*
		pr := &WOFLoggerProgress{
			logger: i.Logger,
		}
	*/

	// something something something auth-y bits
	// https://godoc.org/gopkg.in/src-d/go-git.v4#CloneOptions

	opts := &gogit.CloneOptions{
		URL:   uri,
		Depth: 1,
		// Progress: pr,
	}

	_, err = gogit.PlainClone(tempdir, false, opts)

	if err != nil {
		return err
	}

	return d.repo_driver.IndexURI(ctx, index_cb, tempdir)
}
