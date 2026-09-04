package worker

import (
	"context"
)

type EventType int
const (
	Started EventType = iota
	Progress
	Finished
	Failed
)

type Handler func(context.Context, chan <- Event) error

type Job struct {
	ID    int
	Name string
	Run  func(ctx context.Context, events chan <- Event) error
}

type Event struct {
	JobID   int
	Type    EventType
	Message string
	Err     error
}

type Worker struct {
	Jobs   chan Job
	Events chan Event
	nextID int
	cancel context.CancelFunc
}

func NewWorker() *Worker {
	jobs := make(chan Job)
	events := make(chan Event)
	w := &Worker{
		Jobs: jobs,
		Events: events,
	}
	go w.run()
	return w
}

func (w *Worker) NewJob(name string, fn Handler) Job {
	id := w.nextID
	w.nextID += 1
	return Job{
		ID: id,
		Name: name,
		Run: fn,
	}
}


func (w *Worker) run() {
	
	result := make(chan error, 1)

	for job := range w.Jobs {

		if w.cancel != nil {
			w.cancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel

		w.Events <- Event{
			JobID: job.ID,
			Type:  Started,
			Message: "job started",
		}


		go func() {
			err := job.Run(ctx, w.Events)
			result <- err
		}()

		go func () {
			select {
			case <- ctx.Done():
				w.Events <- Event{
					JobID: job.ID,
					Type:  Failed,
					Err: ctx.Err(),
					Message: "job canceled",
				}
			case err := <- result:
				if err != nil {
					w.Events <- Event{
						JobID:   job.ID,
						Type:    Failed,
						Err:     err,
						Message: "job failed",
					}
					return
				}
				w.Events <- Event{
					JobID: job.ID,
					Type: Finished,
					Message: "job finished",
				}
			}
		}()
	}
}

func NewFailedEvent(err error, jobID int) Event {
	return Event{
		JobID: jobID,
		Type: Failed,
		Message: "job failed",
		Err: err,
	}
}

