package app

import (
	"testing"

	"github.com/dubbersthehoser/mayble/internal/command"
)


func TestManager(t *testing.T) {
	manager := newManager()

	cmd := &command.StubCommand{
		Label: "command-1",
	}

	// Check Execution
	err := manager.execute(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if cmd.Count != 1 {
		t.Fatalf("expect count %d, got %d", 1, cmd.Count)
	}
	if manager.undos.Length() != 1 {
		t.Fatalf("expect length %d, got %d", 1, manager.undos.Length())
	}
	if manager.redos.Length() != 0 {
		t.Fatalf("expect length %d, got %d", 0, manager.redos.Length())
	}

	// Check Un-Execution
	err = manager.unExecute()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if cmd.Count != 0 {
		t.Fatalf("expect count %d, got %d", 0, cmd.Count)
	}
	if manager.undos.Length() != 0 {
		t.Fatalf("expect length %d, got %d", 0, manager.undos.Length())
	}
	if manager.redos.Length() != 1 {
		t.Fatalf("expect length %d, got %d", 1, manager.redos.Length())
	}

	// Check Re-Execution
	err = manager.reExecute()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if cmd.Count != 1 {
		t.Fatalf("expect count %d, got %d", 1, cmd.Count)
	}
	if manager.undos.Length() != 1 {
		t.Fatalf("expect length %d, got %d", 1, manager.undos.Length())
	}
	if manager.redos.Length() != 0 {
		t.Fatalf("expect length %d, got %d", 0, manager.redos.Length())
	}

	// Check Execute with Redos
	err = manager.unExecute()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	cmd = &command.StubCommand{
		Label: "command-2",
	}
	err = manager.execute(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if manager.redos.Length() != 0 {
		t.Fatalf("expect length %d, got %d", 0, manager.redos.Length() )
	}

}








