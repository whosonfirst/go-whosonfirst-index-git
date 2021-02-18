package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/whosonfirst/go-whosonfirst-index-git/v2"
	"github.com/whosonfirst/go-whosonfirst-index/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-index/v2/indexer"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

func main() {

	valid_schemes := strings.Join(emitter.Schemes(), ",")
	emitter_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-index/v2/emitter URI. Supported emitter URI schemes are: %s", valid_schemes)

	var emitter_uri = flag.String("emitter-uri", "repo://", emitter_desc)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

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
