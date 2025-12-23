package command

import (
	"testing"
)

func TestStack(t *testing.T) {

	stack := NewStack()

	cmdStub := &StubCommand{Label: "command-1"}
	stack.Push(cmdStub)
	cmdStub = &StubCommand{Label: "command-2"}
	stack.Push(cmdStub)
	cmdStub = &StubCommand{Label: "command-3"}
	stack.Push(cmdStub)

	if stack.Length() != 3 {
		t.Fatalf("want length %d, got %d", 3, stack.Length())
	}

	cmd := stack.Pop().(*StubCommand)
	if cmd.Label != "command-3" {
		t.Fatalf("want %s, got %s", "command-3", cmd.Label)
	}
	cmd = stack.Pop().(*StubCommand)
	if cmd.Label != "command-2" {
		t.Fatalf("want %s, got %s", "command-2", cmd.Label)

	}
	cmd = stack.Pop().(*StubCommand)
	if cmd.Label != "command-1" {
		t.Fatalf("want %s, got %s", "command-1", cmd.Label)
	}

	last := stack.Pop()
	if last != nil {
		t.Fatalf("want %v, got %v", nil, last)
	}

	for _ = range 10 {
		stack.Push(cmdStub)
	}

	if stack.Length() != 10 {
		t.Fatalf("want %d, got %d", 10, stack.Length())
	}

	stack.Clear()
	if stack.Length() != 0 {
		t.Fatalf("want %d, got %d", 0, stack.Length())
	}
}

func TestQueue(t *testing.T) {
	queue := NewQueue()

	cmdStub := &StubCommand{Label: "command-1"}
	queue.Enqueue(cmdStub)
	cmdStub = &StubCommand{Label: "command-2"}
	queue.Enqueue(cmdStub)
	cmdStub = &StubCommand{Label: "command-3"}
	queue.Enqueue(cmdStub)

	if queue.Length() != 3 {
		t.Fatalf("want length %d, got %d", 3, queue.Length())
	}

	cmd := queue.Dequeue().(*StubCommand)
	if cmd.Label != "command-1" {
		t.Fatalf("want %s, got %s", "command-1", cmd.Label)
	}
	cmd = queue.Dequeue().(*StubCommand)
	if cmd.Label != "command-2" {
		t.Fatalf("want %s, got %s", "command-2", cmd.Label)
	}
	cmd = queue.Dequeue().(*StubCommand)
	if cmd.Label != "command-3" {
		t.Fatalf("want %s, got %s", "command-3", cmd.Label)
	}

	last := queue.Dequeue()
	if last != nil {
		t.Fatalf("want %v, got %v", nil, last)
	}

	for _ = range 10 {
		queue.Enqueue(cmdStub)
	}

	queue.Clear()
	if queue.Length() != 0 {
		t.Fatalf("want length %d, got %d", 0, queue.Length())
	}
	
}






