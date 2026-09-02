package table

import (
	"github.com/dubbersthehoser/mayble/internal/snapshot"
)

type EventSelected struct {
	Has     bool
	Point   Point
	Version int64
}

type EventSnapshotSelected struct {
	Has      bool
	Point    snapshot.Point
	Versiont int64
}

type EventSorted struct {
	snapshot *snapshot.Snapshot
	ids      []int64
	column   string
	asc      bool
}

type EventColumnHidden struct {
	hidden  []bool
}

type EventSnapshotSearched struct {
	snapshot *snapshot.Snapshot
	searched []snapshot.Point
}

type EventNewSnapshot struct {
	snapshot *snapshot.Snapshot
}

