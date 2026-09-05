package table

import (
	"context"

	"github.com/dubbersthehoser/mayble/internal/worker"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

const (
	JobSearchSnapshot string = "job searching"
	JobLoadingSnapshot string = "job loading"
	JobSorting string = "job searching"
)


func NewJobSearchSnapshot(w *worker.Worker, pattern string, column string) worker.Job {
	job := w.NewJob(JobSearchSnapshot, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		ss := snapshot.Current.Load()
		points, score, err := snapshotSearch(ss, pattern, column)
		if err := ctx.Err(); err != nil {
			return
		}
		if err != nil {
			events <- worker.NewFailedEvent(job.ID, err)
		}
		data := EventSnapshotSearched{
			Version: ss.Version(),
			Points: points,
			Scores: score,
		}
		events <- worker.NewFinished(job.ID, )
	}
	return job
}
