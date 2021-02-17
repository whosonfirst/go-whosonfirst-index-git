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

	var dsn = flag.String("dsn", "git://", "")
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
	
	i, err := index.NewIndexer(ctx, *dsn, cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	err = i.Index(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(count, i.Indexed)
}
