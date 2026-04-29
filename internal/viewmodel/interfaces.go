package viewmodel

type TableConfig interface {
	SetColumnWidth(header string, width float32)
	GetColumnWidth(header string) (width float32)
	SetHiddenColumns(header []string)
	GetHiddenColumns() []string
}

type UIConfig interface {
	SetWindowBody(int)
	GetWindowBody() int
	TableConfig
}

type DatabaseOpener interface {
	OpenDB(string) error
}
