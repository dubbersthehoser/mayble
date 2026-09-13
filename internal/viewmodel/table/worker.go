package table

import (
	"context"

	"github.com/dubbersthehoser/mayble/internal/worker"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
	"github.com/dubbersthehoser/mayble/internal/event"
)

const (
	JobSearchTable  string = "searching table"
	JobSortTable    string = "sorting table"
)

func NewJobSearchTable(w *worker.Worker, pattern string, column string) worker.Job {
	job := w.NewJob(JobSearchTable, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		defer close(events)
		ss := snapshot.Current.Load()
		trv, err := getSnapshotTraverser(ss, column)
		if err != nil {
			events <- worker.NewFailedEvent(job.Name, job.ID, event.TableSearched{}, err)
			return
		}
		srch := (&search.Searcher{}).Set(trv, pattern)
		results := searchSearcherWithContext(ctx, srch)

		scores := make([]int, len(results))
		points := make([]models.Cell, len(results))
		for i, r := range results {
			if ctx.Err() != nil {
				return
			}
			scores[i] = r.Score
			id, _ := ss.RowToID(r.Point.Row)
			label := models.BookEntryFields()[r.Point.Col]
			points[i].Column = label
			points[i].ID = id
		}
		data := event.TableSearched{
			Version: ss.Version(),
			Points:  points,
			Scores:  scores,
		}
		events <- worker.NewFinishedEvent(job.Name, job.ID, data)
	}
	return job
}

func NewJobSortTable(w *worker.Worker, column string, asc bool) worker.Job {
	job := w.NewJob(JobSortTable, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		defer close(events)
		ss := snapshot.Current.Load()
		sorted, err := snapshotSort(ss, column, asc)
		if err != nil {
			events <- worker.NewFailedEvent(job.Name, job.ID, event.TableSorted{}, err)
			return
		}
		data := event.TableSorted{
			Version: ss.Version(),
			Sorted:  sorted,
			Asc:     asc,
			Column:  column,
		}
		if ctx.Err() != nil {
			return
		}
		events <- worker.NewFinishedEvent(job.Name, job.ID, data)
	}
	return job
}
