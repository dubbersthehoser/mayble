package gui


const (
	EventSave string = "ui.save"
	EventSaveDisable = "ui.save.disable"
	EventSaveEnable  = "ui.save.enable"

	EventRedo        = "ui.redo"
	EventRedoEmpty   = "ui.redo.empty"
	EventRedoReady   = "ui.redo.ready"

	EventUndo      = "ui.undo"
	EventUndoEmpty = "ui.undo.empty"
	EventUndoReady = "ui.undo.ready"

	EventEditOpen    = "ui.editor.open"
	EventEntryCreate = "ui.entry.create"
	EventEntryDelete = "ui.entry.delete"
	EventEntryUpdate = "ui.entry.update"

	EventEntrySelected   = "ui.entry.selected"
	EventEntryUnselected = "ui.entry.unselected"

	EventSelectNext = "ui.select.next"
	EventSelectPrev = "ui.select.previous"

	EventListOrderBy    = "ui.list.orderby"
	EventListOrdering    = "ui.list.ordering"

	EventSearchBy      = "ui.search.by"
	EventSearchPattern = "ui.search.pattern"



	EventDisplayErr    = "ui.display.error"
)
