package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index/v2/indexer"
	"github.com/whosonfirst/go-whosonfirst-index/v2/emitter"	
	_ "github.com/whosonfirst/go-whosonfirst-index-git/v2"
	"io"
	"log"
	"sync/atomic"
)

func main() {

	var emitter_uri = flag.String("emitter-uri", "git://", "")
	flag.Parse()

	var count int64
	count = 0

	cb := func(ctx context.Context, fh io.ReadSeekCloser, args ...interface{}) error {

		_, err := emitter.PathForContext(ctx)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count, 1)
		return nil
	}

	ctx := context.Background()
	
	idx, err := indexer.NewIndexer(ctx, *emitter_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	err = idx.Index(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Counted %d records (indexed %d records)\n", count, idx.Indexed)
}
