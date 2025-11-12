package memdb

import (
	"github.com/dubbersthehoser/mayble/data"
	"github.com/dubbersthehoser/mayble/storage"
)

type MemStorage struct {
	Books map[int64]data.Book
	Loans map[int64]data.Loan
}

func NewMemStorage() *MemStorage {
	mem := &MemStorage{
		Books: make(map[int64]data.Book),
		Loans:  make(map[int64]data.Loan),
	}
	return mem
}

func (m *MemStorage) GetNewBookID() int64 {
	return int64(len(m.Books)) + 1
}

func (m *MemStorage) GetNewLoanID() int64 {
	return int64(len(m.Loans)) + 1
}

func (m *MemStorage) GetAllBookLoans() ([]data.BookLoan, error) {
	booksLoans := make([]data.BookLoan, len(m.Books))
	count := 0
	for id := range m.Books {
		book := m.Books[id]
		bookLoan := data.BookLoan{
			Book: book,
		}
		loan, ok := m.Loans[id]
		if ok {
			bookLoan.Loan = &loan
		}
		booksLoans[count] = bookLoan
		count++
	}
	return booksLoans, nil
}

func (m *MemStorage) GetBookLoanByID(id int64) (data.BookLoan, error) {
	book, ok := m.Books[id]
	if !ok {
		return data.BookLoan{}, storage.ErrEntryNotFound
	}
	bookLoan := data.BookLoan{
		Book: book,
	}
	loan, ok := m.Loans[id]
	if ok {
		bookLoan.Loan = &loan
	}
	return bookLoan, nil
}

func (m *MemStorage) CreateBookLoan(book *data.BookLoan) (int64, error) {
	if book.ID == data.ZeroID {
		book.ID = m.GetNewBookID()
	}
	_, ok := m.Books[book.ID]
	if ok {
		return 0, storage.ErrEntryExists
	}
	m.Books[book.ID] = book.Book
	if book.IsOnLoan() {
		loanID := m.GetNewLoanID()
		book.Loan.ID = loanID
		m.Loans[book.ID] = *book.Loan
	}
	return book.ID, nil
}

func (m *MemStorage) UpdateBookLoan(book *data.BookLoan) error {
	_, ok := m.Books[book.ID]
	if !ok {
		return storage.ErrEntryNotFound
	}

	m.Books[book.ID] = book.Book

	_, ok = m.Loans[book.ID]

	if !book.IsOnLoan() && ok {
		delete(m.Loans, book.ID)
		return nil
	}

	if book.IsOnLoan() {
		m.Loans[book.ID] = *book.Loan
	}
	return nil

}

func (m *MemStorage) DeleteBookLoan(book *data.BookLoan) error {
	_, ok := m.Books[book.ID]
	if !ok {
		return storage.ErrEntryNotFound
	}
	delete(m.Books, book.ID)
	_, ok = m.Loans[book.ID]
	if ok {
		delete(m.Loans, book.ID)
	}
	return nil
}

func (m *MemStorage) Close() error {
	return nil
}


