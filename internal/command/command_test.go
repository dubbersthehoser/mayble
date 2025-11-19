package command

import (
	"testing"

	"github.com/dubbersthehoser/mayble/internal/command/stub"
)


func TestStack(t *testing.T) {

	stack := NewStack()

	cmdStub := &stub.Command{Label: "command-1"}
	stack.Push(cmdStub)
	cmdStub = &stub.Command{Label: "command-2"}
	stack.Push(cmdStub)
	cmdStub = &stub.Command{Label: "command-3"}
	stack.Push(cmdStub)

	if stack.Length() != 3 {
		t.Fatalf("want %d, got %d", 3, stack.Length())
	}

	cmd := stack.Pop().(*stub.Command)
	if cmd.Label != "command-3" {
		t.Fatalf("want %s, got %s", "command-3", cmd.Label)
	}
	cmd = stack.Pop().(*stub.Command)
	if cmd.Label != "command-2" {
		t.Fatalf("want %s, got %s", "command-2", cmd.Label)

	}
	cmd = stack.Pop().(*stub.Command)
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
