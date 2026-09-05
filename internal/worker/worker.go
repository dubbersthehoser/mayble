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

type Handler func(context.Context, chan <- Event)

type Job struct {
	ID    int
	Name string
	Run  Handler
}

type Event struct {
	JobID   int
	Type    EventType
	Message string
	Err     error
	Data    any
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
			job.Run(ctx, w.Events)
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
			}
		}()
	}
}

func NewFailedEvent(jobID int, err error) Event{
	return Event{
		JobID: jobID,
		Type: Failed,
		Message: "job failed",
		Err: err,
	}
}

func NewFinishedEvent(jobID int, data any) Event{
	return Event{
		JobID: jobID,
		Type: Finished,
		Message: "job finished",
		Err: nil,
		Data: data,
	}
}
