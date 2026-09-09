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
	JobName string
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
			<- ctx.Done()
			w.Events <- Event{
				JobID: job.ID,
				Type:  Failed,
				Err: ctx.Err(),
				Message: "job canceled",
			}
		}()
	}
}

func NewFailedEvent(name string, jobID int, data any, err error) Event{
	return Event{
		JobID: jobID,
		JobName: name,
		Type: Failed,
		Message: "job failed",
		Data: data,
		Err: err,
	}
}

func NewFinishedEvent(name string, jobID int, data any) Event{
	return Event{
		JobID: jobID,
		JobName: name,
		Type: Finished,
		Message: "job finished",
		Err: nil,
		Data: data,
	}
}
