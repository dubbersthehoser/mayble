package data

type Book struct {
	ID      int64  
	Title   string 
	Author  string 
	Genre   string 
	Ratting int    
}

type Loan struct {
	ID   int64     
	Name string   
	Date time.Time 
}

type BookLoan struct {
	Book
	Loan *Loan
}

func NewBookLoan() *BookLoan {
	b := &BookLoan{
		Book: Book{
			ID: ZeroID,
		},
		Loan: &Loan{
			ID: ZeroID,
		},
	}
	return b
}

func (bl *BookLoan) IsOnLoan() bool {
	return bl.Loan != nil
}
func (bl *BookLoan) UnsetLoan() {
	bl.Loan = nil
}
