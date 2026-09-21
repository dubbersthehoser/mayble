package table

import (
	"context"
	"testing"
	"time"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

func Test_tableTraverse(t *testing.T) {

	books := []models.BookEntry{}

	pattern := "author"
	ss := snapshot.NewSnapshot(books)
	trv := newTableTraverse(ss)
	srch := (&search.Searcher{}).Set(trv, pattern)

	completed := make(chan struct{})
	outoftime := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_ = searchSearcherWithContext(ctx, srch)
		completed <- struct{}{}
		close(completed)
	}()

	go func() {
		time.Sleep(time.Millisecond)
		outoftime <- struct{}{}
		cancel()
		close(outoftime)
	}()

	select {
	case <-completed:
	case <-outoftime:
		t.Fatalf("inf loop detected.")
	}
}
