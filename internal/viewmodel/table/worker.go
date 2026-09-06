package table

import (
	"context"

	"github.com/dubbersthehoser/mayble/internal/worker"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

const (
	JobSearchSnapshot  string = "job searching"
	JobLoadingSnapshot string = "job loading"
	JobSortSnapshot    string = "job searching"
)

func NewJobSearchSnapshot(w *worker.Worker, pattern string, column string) worker.Job {
	job := w.NewJob(JobSearchSnapshot, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		ss := snapshot.Current.Load()
		points, score, err := snapshotSearchWithContext(ctx, ss, pattern, column)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			events <- worker.NewFailedEvent(job.Name, job.ID, err)
			return
		}
		data := EventSnapshotSearched{
			Version: ss.Version(),
			Points: points,
			Scores: score,
		}
		events <- worker.NewFinishedEvent(job.Name, job.ID, data)
	}
	return job
}

func NewJobSortSnapshot(w *worker.Worker, column string, asc bool) worker.Job {
	job := w.NewJob(JobSortSnapshot, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		ss := snapshot.Current.Load()
		sorted, err := snapshotSort(ss, column, asc)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			events <- worker.NewFailedEvent(job.Name, job.ID, err)
			return
		}
		data := EventSnapshotSorted{
			Version: ss.Version(),
			Sorted: sorted,
			Asc: asc,
			Column: column,
		}
		events <- worker.NewFinishedEvent(job.Name, job.ID, data)
	}
	return job
}
