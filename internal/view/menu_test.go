package view

import (
	"testing"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func TestTableMenu(t *testing.T) {

	const (
		CheckedOnHidden bool = false
		CheckedOnShow   bool = true
	)

	tests := []struct {
		Label     string
		setHidden func(*ColumnHiddenChanged, bool)
	}{
		{
			Label: "Show Loaned",
			setHidden: func(ev *ColumnHiddenChanged, t bool) {
				ev.LoanSet = t
			},
		},
		{
			Label: "Show Read",
			setHidden: func(ev *ColumnHiddenChanged, t bool) {
				ev.ReadSet = t
			},
		},
		{
			Label: "Show ID",
			setHidden: func(ev *ColumnHiddenChanged, t bool) {
				ev.ID = t
			},
		},
	}

	cb := NewCommandBus()
	eb := NewEventBus()
	menu := newTableMenu(eb, cb)

	if len(menu.Items) != len(tests) {
		t.Fatalf("expect %d, got %d", len(tests), len(menu.Items))
	}

	for i, tt := range tests {
		t.Run(tt.Label, func(t *testing.T) {
			item := menu.Items[i]
			if tt.Label != item.Label {
				t.Fatalf("expect %s, got %s", tt.Label, item.Label)
			}
			ev := ColumnHiddenChanged{}
			tt.setHidden(&ev, false)
			eb.Notify(ev)
			if item.Checked != CheckedOnShow {
				t.Fatalf("expect %t, got %t", CheckedOnShow, item.Checked)
			}
			tt.setHidden(&ev, true)
			eb.Notify(ev)
			if item.Checked != CheckedOnHidden {
				t.Fatalf("expect %t, got %t", CheckedOnHidden, item.Checked)
			}

			// disable button based on what body is changed.
			eb.Notify(BodyChanged{Body: viewmodel.BodyNoData})
			if item.Disabled != true {
				t.Fatalf("expect %t, got %t", true, item.Disabled)
			}
			eb.Notify(BodyChanged{Body: viewmodel.BodyTable})
			if item.Disabled != false {
				t.Fatalf("expect %t, got %t", false, item.Disabled)
			}

			var toggle bool
			cb.Regester("ToggleHiddenColumn", func(v any) error {
				toggle = !toggle
				return nil
			})
			item.Action()
			if !toggle {
				t.Fatalf("dispatch was not called")
			}
			item.Action()
			if toggle {
				t.Fatalf("dispatch did not toggle value")
			}
		})
	}

}
