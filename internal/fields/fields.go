package fields

const Length int = 7

const (
	Title  int = iota
	Author
	Genre
	
	Read
	Rating

	Loaned
	Borrower
)

func Headers() []string {
	return []string{
		"Title",
		"Author",
		"Genre",
		"Read",
		"Rating",
		"Loaned",
		"Borrower",
	}
}
