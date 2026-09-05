package worker

import (
	"testing"
	"context"
	"time"
)

func TestWorker(t *testing.T) {
	
	worker := NewWorker()

	totalTime := time.Second * 2

	run := func(ctx context.Context, ch chan <- Event) error {
			t.Log("Running")
			time.Sleep(totalTime)
			return nil
		}

	job := worker.NewJob("test-one", run)

	worker.Jobs <- job

	go func() {
		job := worker.NewJob("test-two", run)
		time.Sleep(totalTime / 2)
		t.Logf("passing job: %d", job.ID)
		worker.Jobs <- job
	}()

	finished := 0
	canceled := 0

	for event := range worker.Events {
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
