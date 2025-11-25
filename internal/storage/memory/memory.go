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


func (m *Storage) Wipe() {
	for id := range m.Books {
		delete(m.Books, id)
	}
	for id := range m.Loans {
		delete(m.Loans, id)
	}
}


/************************
        StoreBook
*************************/

func (m *Storage) CreateBook(id int64, title, author, genre string, ratting int) (int64, error) {
	if id < 0 {
		id = m.getNewBookID()
	}
	_, ok := m.Books[id]
	if ok {
		return -1, storage.ErrEntryExists
	}
	book := storage.Book{
		ID: id,
		Title: title,
		Author: author,
		Genre: genre,
		Ratting: ratting,
	}
	m.Books[id] = book
	return id, nil
}

func (m *Storage) UpdateBook(book *storage.Book) error {
	if book == nil || book.ID < 0 {
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
	if book == nil || book.ID < 0 {
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

func (m *Storage) GetBookByID(id int64) (*storage.Book, error) {
	if id < 0 {
		return nil, storage.ErrInvalidValue
	}
	book, ok := m.Books[id]
	if !ok {
		return nil, storage.ErrEntryNotFound
	}
	return &book, nil

}


/************************
        StoreLoan
*************************/

func (m *Storage) CreateLoan(ID int64, borrower string, date time.Time) error {
	if ID < 0 {
		return storage.ErrInvalidValue
	}
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
	if loan == nil || loan.ID < 0 {
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
	if loan == nil || loan.ID < 0 {
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




