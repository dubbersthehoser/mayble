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
	Jobs  chan <- Job
	Events <- chan Event

	jobs  <- chan Job
	events chan <- Event
}

func NewWorker() *Worker {
	jobs := make(chan Job)
	events := make(chan Event)
	w := &Worker{
		jobs: jobs,
		Jobs: jobs,
		Events: events,
		events: events,
	}
	go w.run()
	return w
}


func (w *Worker) run() {
	for job := range w.jobs {
		ctx := context.Background()

		w.events <- Event{
			JobID: job.ID,
			Type:  Started,
			Message: "job started",
		}

		var event Event
		event = Event{
			JobID: job.ID,
			Type: Finished,
			Message: "job finished",
		}

		err := job.Run(ctx, w.events)
		if err != nil {
			event = Event{
				JobID: job.ID,
				Type: Failed,
				Message: "job failed",
			}
		}
		w.events <- event
	}
}


