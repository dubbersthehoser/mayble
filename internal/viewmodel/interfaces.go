package viewmodel

type TableConfigorator interface {
	SetColumnWidth(header string, width float32)
	GetColumnWidth(header string) (width float32)
	SetHiddenColumns(header []string)
	GetHiddenColumns() []string
}

type UIConfig interface {
	SetWindowBody(int)
	GetWindowBody() int
	TableConfigorator
}

type DatabaseOpener interface {
	OpenDB(string) error
}
