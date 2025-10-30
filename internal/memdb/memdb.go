package memdb

import (
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type MemStorage struct {
	Books map[int64]storage.Book
	Loans map[int64]storage.Loan
}

func NewMemStorage() *MemStorage {
	mem := &MemStorage{
		Books: make(map[int64]storage.Book),
		Loans:  make(map[int64]storage.Loan),
	}
	return mem
}

func (m *MemStorage) GetNewBookID() int64 {
	return int64(len(m.Books)) + 1
}

func (m *MemStorage) GetNewLoanID() int64 {
	return int64(len(m.Loans)) + 1
}

func (m *MemStorage) GetAllBookLoans() ([]storage.BookLoan, error) {
	booksLoans := make([]storage.BookLoan, len(m.Books))
	count := 0
	for id := range m.Books {
		book := m.Books[id]
		bookLoan := storage.BookLoan{
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

func (m *MemStorage) GetBookLoanByID(id int64) (storage.BookLoan, error) {
	book, ok := m.Books[id]
	if !ok {
		return storage.BookLoan{}, storage.ErrEntryNotFound
	}
	bookLoan := storage.BookLoan{
		Book: book,
	}
	loan, ok := m.Loans[id]
	if ok {
		bookLoan.Loan = &loan
	}
	return bookLoan, nil
}

func (m *MemStorage) CreateBookLoan(book *storage.BookLoan) (int64, error) {
	if book.ID == storage.ZeroID {
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

func (m *MemStorage) UpdateBookLoan(book *storage.BookLoan) error {
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

func (m *MemStorage) DeleteBookLoan(book *storage.BookLoan) error {
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


