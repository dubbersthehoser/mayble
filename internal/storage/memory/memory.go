package memory

import (
	"time"

	"github.com/dubbersthehoser/mayble/internal/storage"
)

type Storage struct {
	Books map[int64]storage.Book
	Loans map[int64]storage.Loan
}

func NewStorage() *Storage {
	mem := &Storage{
		Books: make(map[int64]storage.Book),
		Loans:  make(map[int64]storage.Loan),
	}
	return mem
}

/************************
        StoreBook
*************************/

func (m *Storage) CreateBook(title, author, genre string, ratting int) (int64, error) {
	book := storage.Book{
		Title: title,
		Author: author,
		Genre: genre,
		Ratting: ratting,
	}
	id := m.getNewBookID()
	m.Books[id] = book
	return id, nil
}

func (m *Storage) UpdateBook(book *storage.Book) error {
	if book == nil {
		return storage.ErrInvalidValue
	}
	_, ok := m.Books[book.ID]
	if !ok {
		return storage.ErrEntryNotFound
	}
	m.Books[book.ID] = *book
	return nil
}

func (m *Storage) DeleteBook(book *storage.Book) error {
	if book == nil {
		return storage.ErrInvalidValue
	}
	_, ok := m.Books[book.ID]
	if !ok {
		return storage.ErrEntryNotFound
	}
	delete(m.Books, book.ID)
	return nil
}

func (m *Storage) GetBooks() ([]storage.Book, error) {
	r := make([]storage.Book, len(m.Books))
	i := 0
	for _, v := range m.Books {
		r[i] = v
		i+=1
	}
	return r, nil
}


/************************
        StoreLoan
*************************/

func (m *Storage) CreateLoan(ID int64, borrower string, date time.Time) error {
	_, ok := m.Loans[ID]
	if ok {
		return storage.ErrEntryExists
	}
	loan := storage.Loan{
		ID: ID,
		Borrower: borrower,
		Date: date,
	}
	m.Loans[ID] = loan
	return nil
}

func (m *Storage) UpdateLoan(loan *storage.Loan) error {
	if loan == nil {
		return storage.ErrInvalidValue
	}
	_, ok := m.Loans[loan.ID]
	if !ok {
		return storage.ErrEntryNotFound
	}
	m.Loans[loan.ID] = *loan
	return nil
}

func (m *Storage) DeleteLoan(loan *storage.Loan) error {
	if loan == nil {
		return storage.ErrInvalidValue
	}
	_, ok := m.Loans[loan.ID]
	if !ok {
		return storage.ErrEntryNotFound
	}
	delete(m.Loans, loan.ID)
	return nil
}

func (m *Storage) GetLoan(id int64) (*storage.Loan, error) {
	loan, ok := m.Loans[id]
	if !ok {
		return &loan, storage.ErrEntryNotFound
	}
	return &loan, nil
}

func (m *Storage) getNewBookID() int64 {
	return int64(len(m.Books)) + 1
}




