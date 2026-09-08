package table

import (
	"testing"

	"github.com/dubbersthehoser/mayble/internal/models"
)

func TestRoundTrip_toColumn(t *testing.T) {
	tests := []struct{
		name        string
		header      []string
		sheetCol    int
		snapshotCol int
		willErr     bool
	}{
		{
			name: "test-0: with only Title and Author",
			header: []string{
				models.BookEntryFields()[models.IdxTitle],
				models.BookEntryFields()[models.IdxAuthor],
			},
			sheetCol: 0,
			snapshotCol: models.IdxTitle,
			willErr: false,
		},
		{
			name: "test-1: with only Title and Completed",
			header: []string{
				models.BookEntryFields()[models.IdxTitle],
				models.BookEntryFields()[models.IdxCompletedAt],
			},
			sheetCol: 1,
			snapshotCol: models.IdxCompletedAt,
			willErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := toSheetColumn(tt.header, tt.snapshotCol)
			if tt.willErr {
				if err == nil {
					t.Fatal("expected err")
				}
				return
			}
			if actual != tt.sheetCol {
				t.Fatalf("expect %d, got %d", tt.sheetCol, actual)
			}

			actual, err = toSnapshotColumn(tt.header, tt.sheetCol)
			if tt.willErr {
				if err == nil {
					t.Fatal("expected err")
				}
				return
			}
			if actual != tt.snapshotCol {
				t.Fatalf("expect %d, got %d", tt.snapshotCol, actual)
			}
		})
	}
}

