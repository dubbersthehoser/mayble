package viewmodel

type UIConfig interface {
	SetColumnWidth(header string, width float32)
	GetColumnWidth(header string) (width float32)
	SetHiddenColumns(header []string)
	GetHiddenColumns() []string
}
