package git

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	
)

func init() {
	dr := NewGitDriver()
	index.Register("git", dr)
}

type GitDriver struct {
	index.Driver
}

func NewGitDriver() index.Driver {

	dr := &GitDriver{}

	return dr
}

func (d *GitDriver) Open(uri string) error {

	return nil
}

func (d *GitDriver) IndexURI(ctx context.Context, index_cb index.IndexerFunc, uri string) error {

	opts := &gogit.CloneOptions{
		URL: uri,
	}
	
	r, err := gogit.Clone(memory.NewStorage(), nil, opts)

	if err != nil {
		return err
	}

	it, err := r.BlobObjects()

	if err != nil {
		return err
	}

	err = it.ForEach(func(bl *object.Blob) error {

		fh, err := bl.Reader()

		if err != nil {
			return err
		}

		defer fh.Close()

		return index_cb(ctx, fh)
	})

	return nil
}
