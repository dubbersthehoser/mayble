package app

import (
	"context"

	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/worker"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

const (
	JobImportFile     string = "importing file"
	JobExportDatabase string = "exporting database"
	JobTakeSnapshot   string = "loading snapshot"
)

func NewJobTakeSnapshot(w *worker.Worker, srv *Service) worker.Job {
	job := w.NewJob(JobTakeSnapshot, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		defer close(events)
		books, err := srv.getAllBooksWithContext(ctx)
		if err != nil {
			events <- worker.NewFailedEvent(job.Name, job.ID, event.StoredSnapshot{
				Message: err.Error(),
				Failed: true,
			}, err)
			return
		}
		ss := snapshot.NewSnapshot(books)
		snapshot.Current.Store(ss)
		events <- worker.NewFinishedEvent(job.Name, job.ID, event.StoredSnapshot{
			Version: ss.Version(),
			Failed: false,
		})
	}
	return job
}

func NewJobImportFile(w *worker.Worker, srv *Service, path string) worker.Job {
	job := w.NewJob(JobImportFile, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		defer close(events)
		err := srv.importFileWithContext(ctx, path)
		var e worker.Event
		if err != nil {
			e = worker.NewFailedEvent(job.Name, job.ID, event.ImportedFile{
				Failed: true,
				Message: err.Error(),
				Path: path,
			}, err)
		} else {
			e =  worker.NewFinishedEvent(job.Name, job.ID, event.ImportedFile{
				Failed: false,
				Message: "imported",
				Path: path,
			})
		}
		events <- e
	}
	return job
}

func NewJobExportFile(w *worker.Worker, srv *Service, path string) worker.Job {
	job := w.NewJob(JobExportDatabase, nil)
	job.Run = func(ctx context.Context, events chan <- worker.Event) {
		defer close(events)
		err := srv.exportFileWithContext(ctx, path)
		var e worker.Event
		if err != nil {
			e = worker.NewFailedEvent(job.Name, job.ID, event.ExportedFile{
				Failed: true,
				Message: "exported",
				Path: path,
			}, err)
		} else {
			e = worker.NewFinishedEvent(job.Name, job.ID, event.ExportedFile{
				Failed: false,
				Path: path,
			})
		}
		events <- e
	}
	return job
}

