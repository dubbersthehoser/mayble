package table

import (
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

type CommandSnapshotSelect struct {
	Version int64
	Point   snapshot.Point
	Has     bool
}

type CommandSheetSelect struct {
	Point Point
	Has   bool
}

type CommandSort struct {
	Column string
	Asc    bool
}

type CommandSearch struct {
	Column  string
	Pattern string
}
