package repository


type BookJoin string

const (
	Main   BookJoin = "Book"
	Read   BookJoin = "BookRead"
	Loaned BookJoin = "BookLoaned"
)


type SortDirection int

const (
	ASC SortDirection = iota
	DESC
)

func (sd SortDirection) String() string {
	switch sd {
	case ASC:
		return "ASC"
	case DESC:
		return "DESC"
	default:
		panic("invalid sort direction")
	}
}


type Resultable interface {
	ID()   int64
	Type() string
}


type ResultSet struct {
	Items []Resultable
	Fields []string
}


type BookSearchParams struct {
	Query     string
	QueryBy   string

	Join      BookJoin

	SortBy    string
	SortOrder SortDirection
}

type BookSearcher interface {
	BookSearch(params BookSearchParams) (ResultSet, error)
}

