package gui


const (
	EventRedo        = "ui.redo"
	EventRedoEmpty   = "ui.redo.empty"
	EventRedoReady   = "ui.redo.ready"

	EventUndo      = "ui.undo"
	EventUndoEmpty = "ui.undo.empty"
	EventUndoReady = "ui.undo.ready"

	EventEditerOpen  = "ui.editor.open"
	EventEntryCreate = "ui.entry.create"
	EventEntryDelete = "ui.entry.delete"
	EventEntryUpdate = "ui.entry.update"
	EventEntrySubmit = "ui.entry.submit"

	EventDocumentModified = "ui.document.modified"
	EventDocumentImport   = "ui.document.import"
	EventDocumentExport   = "ui.document.export"

	EventEntrySelected   = "ui.entry.selected"
	EventEntryUnselected = "ui.entry.unselected"

	EventSelectNext = "ui.select.next"
	EventSelectPrev = "ui.select.previous"

	EventListOrderBy    = "ui.list.orderby"
	EventListOrdering   = "ui.list.ordering"
	EventListOrdered    = "ui.list.ordered"

	EventSearch        = "ui.search"
	EventSearchBy      = "ui.search.by"
	EventSearchPattern = "ui.search.pattern"

	EventSelection     = "ui.selection"
	EventSelectionNone = "ui.selection.none"
	EventSelectionAll  = "ui.selection.all"

	EventDisplayErr    = "ui.display.error"
	EventMenuOpen      = "ui.menu.open"
)
