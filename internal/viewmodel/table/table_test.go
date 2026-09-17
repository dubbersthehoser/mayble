package table

import (
	"testing"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/worker"
)

func TestTable(t *testing.T) {
	eb := event.NewEventBus()
	cb := command.NewCommandBus()
	cfg := config.NewConfigWithDefaults("/tmp/mayble.test.config")
	w := worker.NewWorker()

	table := NewTable(cfg, w, cb, eb)
	testSelected(t, table)
}

func testSelected(t *testing.T, table *Table) {
	t.Helper()

	notifyCalls := 0
	table.Selected.AddListener(func() {
		notifyCalls += 1
	})

	// Check setting selected.
	cell := models.Cell{ID: 1, Column: "Author"}
	table.Selected.Set(0, cell, true)
	
	if notifyCalls != 1 {
		t.Fatalf("expected notify calls %d, got %d", 1, notifyCalls)
	}

	if table.Selected.has != true {
		t.Fatalf("expected has %t, got %t", true, table.Selected.has)
	}

	if table.Selected.selected != cell {
		t.Fatalf("expected cell %v, got %v", cell, table.Selected.selected)
	}

	// Check for unselecting cell.
	cell = models.Cell{}
	table.Selected.Set(0, cell, false)

	if notifyCalls != 2 {
		t.Fatalf("expected notify calls %d got %d", 2, notifyCalls)
	}

	if table.Selected.has != false {
		t.Fatalf("expected has %t, got %t", false, table.Selected.has)
	}

	if table.Selected.selected != cell {
		t.Fatalf("expected cell %v, got %v", cell, table.Selected.selected)
	}
}

