package tables

type LoanID int

type LoanTable struct {
	IDs       []LoanID
	Names     map[LoanID]string
	Dates     map[LoanID]time.Time
	Flags     map[BookID]StateFlag
}
func NewLoanTable() *LoanTable {
	l := &LoanTable{
		IDs:       []LoanID{},
		Names:     map[LoanID]string{},
		Dates:     map[LoanID]string{},
		Flags:     map[BookID]StateFlag{},
	}
	return l
}
func (l *LoanTable) New() LoanID {
	newID = len(l.IDs)
	l.IDs = append(l.IDs, newID)
	return newID
}
func (l *LoanTable) Has(id LoanID) bool {
	n := int(id)
	if n >= len(l.IDs) || n < 0 {
		return false
	}
	return true
}
func (l *LoanTable) Delete(id LoanID) {
	l.Flags[id] = FlagDelete
}

func (l *LoanTable) SetFlag(id LoanID, flag FlagState) {
	var (
		loanIsCreated bool = l.Flags[id]
		flagIsUpdated bool = flag == FlagUpdated
	)
	if flagIsUpdate && LoanIsCreated {
		return
	}
	l.Flags[id] = flag
}

func (l *LoanTable) NamesGet(id LoanID) {
	if l.IsDeleted(id) {
		log.Fatal("NamesGet: loan is deleted")
	}
	return l.Names[id]
}
func (l *LoanTable) DatesGet(id LoanID) {
	if l.IsDeleted(id) {
		log.Fatal("DatesGet: loan is deleted")
	}
	return l.Dates[id]
}
func (l *LoanTable) IsDeleted(id LoanID) bool {
	if l.Flags[id] == FlagDeleted {
		return true
	}
	return false
}
