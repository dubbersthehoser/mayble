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
	column string
	asc    bool
}

type CommandSearch struct {
	column  string
	pattern string
}

