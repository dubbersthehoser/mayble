package worker

import (
	"testing"
	"context"
	"time"
)

func TestWorker(t *testing.T) {
	
	w := NewWorker()

	totalTime := time.Second * 2

	newRun := func(name string, id int, d time.Duration) Handler {
		return func(ctx context.Context, ch chan <- Event) {
			t.Log("Running")
			time.Sleep(d)
			ch <- NewFinishedEvent(name, id, nil)
		}
	}

	job := w.NewJob("job-one", nil)
	job.Run = newRun(job.Name, job.ID, totalTime)

	w.Jobs <- job

	go func() {
		job := w.NewJob("job-two", nil)
		job.Run = newRun(job.Name, job.ID, totalTime / 2)
		time.Sleep(totalTime / 2)
		t.Logf("passing job: %d", job.ID)
		w.Jobs <- job
	}()

	finished := 0
	canceled := 0
	for event := range w.Events {
		t.Log(event.Message)
		switch {
		case event.Message == "job canceled":
			canceled+=1
			continue
		case event.Type == Finished:
			finished += 1
		}
		if finished == 1 {
			break
		}
	}

	if canceled != 1 {
		t.Fatalf("first job was not canceled from incoming second job.")
	}
}
