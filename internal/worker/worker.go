package worker

import (
	"context"
	"fmt"
)

type EventType int

const (
	Started EventType = iota
	Progress
	Finished
	Failed
)

type Handler func(context.Context, chan<- Event)

type Job struct {
	ID   int
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

func (e *Event) Format() string {
	return fmt.Sprintf("job %d: %s %s", e.JobID, e.JobName, e.Message)

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
		Jobs:   jobs,
		Events: events,
	}
	go w.run()
	return w
}

func (w *Worker) NewJob(name string, fn Handler) Job {
	id := w.nextID
	w.nextID += 1
	return Job{
		ID:   id,
		Name: name,
		Run:  fn,
	}
}

func (w *Worker) run() {
	for job := range w.Jobs {

		if w.cancel != nil {
			w.cancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel

		w.Events <- NewStartedEvent(job.Name, job.ID)

		go w.runJob(ctx, job)
	}
}

func (w *Worker) runJob(ctx context.Context, job Job) {

	out := make(chan Event)
	go func() {
		job.Run(ctx, out)
	}()

	hasCanceled := false

	for {
		select {
		case e, ok := <-out:
			if !ok {
				return
			}
			if hasCanceled {
				continue
			}
			select {
			case w.Events <- e:
			case <-ctx.Done():
				hasCanceled = true
				w.Events <- NewCanceledEvent(job.Name, job.ID, ctx.Err())
			}
		case <-ctx.Done():
			hasCanceled = true
			w.Events <- NewCanceledEvent(job.Name, job.ID, ctx.Err())
			for range out {
			}
			return
		}
	}
}

func NewStartedEvent(name string, jobID int) Event {
	return Event{
		JobID:   jobID,
		Type:    Started,
		JobName: name,
		Message: "started",
	}

}

func NewCanceledEvent(name string, jobID int, err error) Event {
	e := NewFailedEvent(name, jobID, nil, err)
	e.Message = "canceled"
	return e
}

func NewFailedEvent(name string, jobID int, data any, err error) Event {
	return Event{
		JobID:   jobID,
		JobName: name,
		Type:    Failed,
		Message: "failed",
		Data:    data,
		Err:     err,
	}
}

func NewFinishedEvent(name string, jobID int, data any) Event {
	return Event{
		JobID:   jobID,
		JobName: name,
		Type:    Finished,
		Message: "finished",
		Err:     nil,
		Data:    data,
	}
}
