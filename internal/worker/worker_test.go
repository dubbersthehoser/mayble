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
			defer close(ch)
			t.Logf("Running: %s", name)
			time.Sleep(d)
			ch <- NewFinishedEvent(name, id, nil)
		}
	}

	firstJob := w.NewJob("job-one", nil)
	firstJob.Run = newRun(firstJob.Name, firstJob.ID, totalTime)

	secondJob := w.NewJob("job-two", nil)
	secondJob.Run = newRun(secondJob.Name, secondJob.ID, totalTime)

	go func() {
		w.Jobs <- firstJob
		time.Sleep(totalTime/2)
		w.Jobs <- secondJob
	}()

	finished := 0
	canceled := 0
	for event := range w.Events {
		t.Log(event.Message)
		switch {
		case event.Message == "canceled":
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
