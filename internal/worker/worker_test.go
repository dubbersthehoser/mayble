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

	job := Job{
		ID: 123,
		Name: "test-one",
		Run: run,
	}

	worker.Jobs <- job

	go func() {
		job2 := job
		job2.ID = 1234
		job2.Name = "test-two"
		time.Sleep(time.Second)
		t.Logf("passing job: %d", job.ID)
		worker.Jobs <- job2
	}()

	for event := range worker.Events {
		t.Log(event.Message)
		if event.Type == Failed {
			break
		}
	}
}
