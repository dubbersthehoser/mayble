package tables

type BookID int

type BookTable struct {
	IDs       []BookID
	Titles   map[BookID]string
	Authors  map[BookID]string
	Genres   map[BookID]string
	Rattings map[BookID]int
	Flag     map[BookID]StateFlag
}
func NewBookTable() *BookTable {
	t := &BookTable{
		IDs:      []BookID{},
		Titles:   map[BookID]string{},
		Authors:  map[BookID]string{},
		Genres:   map[BookID]string{},
		Rattings: map[BookID]string{},
		Flag:     map[BookID]StateFlag{},
	}
	return t
}
func (b *BookTable) New() BookID {
	newID := len(b.IDs)
	l.IDs = append(l.IDs, newID)
	return newID
}
func (b *BookTable) Has(id BookID) bool {
	n := int(id)
	if n >= len(b.IDs) || n < 0 {
		return false
	}
	return true
}
func (b *BookTable) IsDeleted(id BookID) bool {
	if b.Flags[id] == FlagDelete {
		return true
	}
	return false
}
func (b *BookTable) Delete(id BookID) error {
	if b.Has(id) {
		b.Flag[id] = FlagDelete
		return nil
	} 
	return errors.New("BookTable.Delete: book id not found")
}

