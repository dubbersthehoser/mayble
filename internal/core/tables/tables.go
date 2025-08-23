

const NonID int = -127 // used for nil id values

type StateFlag int
const (
	FlagNothing StateFlag = iota
	FlagCreated
	FlagUpdated
	FlagDelete
)

/*
	BookLoanTable
*/

type BookLoanTable struct {
	BookTable  *BookTable
	LoanTable  *LoanTable
	bookToLoan map[BookID]LoanID
	loanToBook map[LoanID]BookID
}
func NewBookLoanTable *BookLoanTable {
	bl := &BookLoanTable{
		BookTable:  NewBookTable(),
		LoanTable:  NewLoanTable(),
		bookToLoan: map[BookID]LoanID{},
		loanToBook: map[LoanID]BookID{},
	}
	return bl
}

/*
	Loan
*/
type LoanCreateParams struct {
	BookID BookID
	Name   string
	Date   *time.Time
}
func (bl *BookLoanTable) LoanCreate(loan LoanCreateParams) (LoanID, error) {
	_, bookLoanFound := bl.bookToLoan[loan.BookID]
	if bookLoanFound {
		return errors.New("LoanCreate: book has loan entry")
	}
	loanID := bl.LoanTable.New()
	bl.loanToBook[loanID] = loan.BookID
	bl.bookToLoan[loan.BookID] = loanID
	bl.LoanTable.StateFlags[loanID] = FlagCreated
	bl.LoanTable.Names[loanID] = loan.Name
	bl.LoanTable.Dates[loanID] = loan.Date
	return loanID
}

type LoanUpdateParams struct {
	LoanID LoanID
	Name   string
	Date   *time.Time
}
func (bl *BookLoanTable) LoanUpdate(loan LoanUpdateParams) error {
	loanID := loan.LoanID
	var LoanDoseNotExist bool = !bl.LoanTable.Has(loanID)
	if LoanDoseNotExist {
		return errors.New("LoanUpdate: loan ID not found")
	}
	bl.LoanTable.SetFlag(LoanID, FlagUpdated)
	bl.LoanTable.Names[loanID] = loan.Name
	bl.LoanTable.Dates[loanID] = loan.Date
}


/*
	Books
*/
type BookCreateParams struct {
	Title   string
	Author  string
	Genre   string
	Ratting int
}
func (bl *BookLoanTable) BookCreate(book BookCreateParams) BookID {
	bookID := bl.BookTable.New()
	bl.BookTable.StateFlags[bookID] = FlagCreated
	bl.BookTable.Titles[bookID] = book.Title
	bl.BookTable.Authors[bookID] = book.Author
	bl.BookTable.Genre[bookID] = book.Genre
	bl.BookTable.Ratting[bookID] = book.Ratting
	return bookID
}

type BookUpdateParams struct {
	BookID  BookID
	Title   string
	Author  string
	Genre   string
	Ratting int
}
func (bl *BookLoanTable) BookUpdate(book BookUpdateParams) error {
	bookID := book.BookID
	var bookDoesNotExists bool = !(bl.BookTable.Has(bookID))
	if bookDoesNotExists {
		return errors.New("BookUpdate: book not found")
	}
	bl.BookTable.Titles[bookID] = book.Title
	bl.BookTable.Authors[bookID] = book.Author
	bl.BookTable.Genre[bookID] = book.Genre
	bl.BookTable.Ratting[bookID] = book.Ratting
	return nil
}

func (bl *BookLoanTable) UnloanBook(id BookID) {
	loanID, ok := bl.bookToLoan[id] 
	if ok {
		delete(bl.bookToLoan[id])
		delete(bl.loanToBook[loanID])
		bl.LoanTable.Delete(loanID)
	}
}

func (bl *BookLoanTable) GetBookLoan(id BookID) (LoanID, error) {
	loanID, loanIsFound := bl.bookToLoan[id]
	if loanIsFound {
		return loanID, nil
	}
	return loanID(NonID), errors.New("GetLoan: book loan not found")
}

func (bl *BookLoanTable) GetBookTitle(id BookID) string {
	return bl.BookTable.TitlesGet(id)
}
func (bl *BookLoanTable) GetBookAuthor(id BookID) string {
	return bl.BookTable.AuthorsGet(id)
}
func (bl *BookLoanTable) GetBookGenre(id BookID) string {
	return bl.BookTable.GenresGet(id)
}
func (bl *BookLoanTable) GetBookRatting(id BookID) int {
	return bl.BookTable.RattingsGet(id)
}

func (bl *BookLoanTable) GetLoanName(id LoanID) string {
	return bl.LoanTable.NamesGet(id)
}
func (bl *BookLoanTable) GetLoanDate(id LoanID) time.Time {
	return bl.LoanTable.DatesGet(id)
}
