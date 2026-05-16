package viewmodel

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/table"
	"github.com/dubbersthehoser/mayble/internal/models"
)

const (
	notifySelect bool   = true
	stubValue    string = "N/A"
)

// The smallest width a column can be.
const MinColWidth float32 = 100.0

type TableBody struct {
	Table       *Table
	Search      *SearchTable
	Controllers *TableControllersVM
}

// 
type HiddenHeaders struct {
	hidden   map[string]bool
	columns []string
}

func NewHiddenHeaders() *HiddenHeaders {
	h := HiddenHeaders{
		hidden: make(map[string]bool),
		columns: models.BookEntryFields(),
	}
	for _, o := range models.BookEntryFields() {
		h.hidden[o] = false
	}
	return &h
}

func (h *HiddenHeaders) SetOptions(options []string) {
	for key := range h.hidden {
		h.hidden[key] = false
	}
	for _, o := range options {
		switch o {
		case "Read":
			h.hidden[h.columns[models.IdxRating]] = true
			h.hidden[h.columns[models.IdxCompletedAt]] = true
		case "Loaned":
			h.hidden[h.columns[models.IdxLoanedAt]] = true
			h.hidden[h.columns[models.IdxBorrower]] = true
		default:
			h.hidden[o] = true
		}
	}
}

func (h *HiddenHeaders) SetHeaders(headers []string) {
	for _, header := range headers {
		h.hidden[header] = true
	}
}

func (h *HiddenHeaders) Options() []string {
	return []string{
		"Title", "Author", "Genre", "Read", "Loaned",
	}
}

func (h *HiddenHeaders) HiddenOptions() []string {
	options := make([]string, 0)
	if h.hidden[h.columns[models.IdxLoanedAt]] &&
	   h.hidden[h.columns[models.IdxBorrower]] {
		options = append(options, "Loaned")
	}
	if h.hidden[h.columns[models.IdxCompletedAt]] &&
	   h.hidden[h.columns[models.IdxRating]] {
		options = append(options, "Read")
	}
	if h.hidden[h.columns[models.IdxTitle]] {
		options = append(options, h.columns[models.IdxTitle])
	}
	if h.hidden[h.columns[models.IdxAuthor]] {
		options = append(options, h.columns[models.IdxAuthor])
	}
	if h.hidden[h.columns[models.IdxGenre]] {
		options = append(options, h.columns[models.IdxGenre])
	}
	return options
}

func (h *HiddenHeaders) Columns() []string {
	columns := make([]string, 0)
	for _, o := range models.BookEntryFields() {
		if h.hidden[o] {
			columns = append(columns, o)
		}
	}
	return  columns
}

type SearchTable struct {
	table  *Table
	header string
}
func NewSearchTable(t *Table) *SearchTable {
	return &SearchTable{
		table: t,
	}
}

func (st *SearchTable) SetSearchColumn(header string) {
	st.header = header
}

func (st *SearchTable) SearchColumnList() []string {
	return []string{
		"All",
		"Title",
		"Author",
		"Genre",
		"Borrower",
	}
}

func (st *SearchTable) Search(s string) {
	
}

type TableHeader struct {
	table *Table
	cfg   TableHeaderConfigurator
}

func NewTableHeader(t *Table, cfg TableHeaderConfigurator) *TableHeader {
	return &TableVisableHeader{
		table: t,
		cfg: cfg,
	}
}

func (vh *TableHeader) SetHidden(headers []string) {
	
}


type Table struct {
	repo   repo.BookRetriever
	bus    *bus.Bus
	cfg    UIConfig

	SortBy    binding.String
	SortOrder binding.String

	HideHeaders *HiddenHeaders

	selector *EntrySelect

	Search struct {
		Text   binding.String
		Header binding.String
	}

	table *table.Table

	l *listener
}

func NewTableVM(b *bus.Bus, cfg UIConfig, r repo.BookRetriever) *Table {
	t := &Table{
		table:  table.NewTable("Main", entryHeaders()),
		repo:   r,
		cfg:    cfg,
		bus:    b,

		SortBy:    binding.NewString(),
		SortOrder: binding.NewString(),

		HideHeaders: NewHiddenHeaders(),

		Search: struct {
			Text   binding.String
			Header binding.String
		}{
			Text:   binding.NewString(),
			Header: binding.NewString(),
		},

		selector: newEntrySelect(r),

		l: &listener{},
	}

	t.HideHeaders.SetHeaders(cfg.GetHiddenColumns())

	t.Search.Text.AddListener(binding.NewDataListener(func() {
		t.selector.unselect(notifySelect)
		t.search()
	}))

	_ = t.SortOrder.Set("ASC")
	_ = t.SortBy.Set(t.table.Headers()[0])
	err := t.load()
	if err != nil {
		log.Println(err)
	}

	b.Subscribe(bus.Handler{
		Name: msgDataChanged,
		Handler: func(e *bus.Event) {
			t.selector.unselect(notifySelect)
			err := t.reload()
			if err != nil {
				log.Println(err)
				return
			}
			t.l.notify()
		},
	})

	return t
}

func (t *Table) search() {
	search, _ := t.Search.Text.Get()
	if search == "" {
		return
	}
	header, _ := t.Search.Header.Get()
	if header == "All" {
		header = ""
	}
	result := table.Search(t.table, search, header)
	if len(result) == 0 {
		return
	}
	r := result[0]
	t.selector.selectID(r.ID, !notifySelect)
	t.selector.selectCell(r.Row, r.Col, notifySelect)
}

func (t *Table) SetSelector(es *EntrySelect) *Table {
	t.selector = es
	return t
}

// SearchOptions a list of searchable options.
func (t *Table) SearchOptions() []string {
	return []string{
		"All",
		"Title",
		"Author",
		"Genre",
		"Borrower",
	}
}

// Selector returns the table's selector.
func (t *Table) Selector() *EntrySelect {
	return t.selector
}

// Sort table using sort bindings.
func (t *Table) Sort() error {
	t.selector.unselect(notifySelect)
	err := t.reload()
	if err != nil {
		return err
	}
	t.l.notify()
	return nil
}

// StoreColumnWidth to the config file if it exists else it will be an nop.
// When width is smaller then MinColWidth, MinColWidth will be used.
func (t *Table) StoreColumnWidth(col int, width float32) {
	if width < MinColWidth {
		width = MinColWidth
	}
	header := t.table.GetHeader(col)
	label := header.Name()
	t.cfg.SetColumnWidth(label, width)
}

// GetColumnWidth from the config file if it exsits, else returns defualt MinColWidth.
func (t *Table) GetColumnWidth(col int) float32 {
	label := t.table.GetHeader(col).Name()
	width := t.cfg.GetColumnWidth(label)
	if width < MinColWidth {
		width = MinColWidth
	}
	return width
}

// SetHiddenHeaders set named option headers to hidden.
func (t *Table) SetHiddenHeaders(options []string) {
	t.HideHeaders.SetOptions(options)
	t.cfg.SetHiddenColumns(t.HideHeaders.Columns())
	t.table.SetHidden(t.HideHeaders.Columns())
	t.l.notify()
}

// HiddenHeaders returns set named options of hidden headers.
func (t *Table) HiddenHeaders() []string {
	return t.HideHeaders.HiddenOptions()
}

// HiddenOptions returns options for hidding columns.
func (t *Table) HiddenOptions() []string {
	return t.HideHeaders.Options()
}


func (t *Table) Headers() []string {
	return t.table.Headers()
}

func (t *Table) Select(row, col int) {
	cell := t.table.GetCell(row, col)
	t.selector.selectID(cell.ID(), !notifySelect)
	t.selector.selectCell(row, col, notifySelect)
}

func (t *Table) Unselect(row, col int) {
	t.selector.unselect(!notifySelect)
}

// reload clear table, then call load.
func (t *Table) reload() error {
	t.table.ClearValues()
	return t.load()
}

// sortbooks sort slice of books.
func sortBooks(books []models.BookEntry, header string, desending bool) error {

	index := slices.Index(entryHeaders(), header)

	if index == -1 {
		return fmt.Errorf("sort_books: invalid header '%s'", header)
	}

	slices.SortFunc(books, func(a, b models.BookEntry) int {
		r := -1
		switch index {
		case models.IdxTitle:
			r = cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
		case models.IdxAuthor:
			r = cmp.Compare(strings.ToLower(a.Author), strings.ToLower(b.Author))
		case models.IdxGenre:
			r = cmp.Compare(strings.ToLower(a.Genre), strings.ToLower(b.Genre))
		case models.IdxBorrower:
			r = cmp.Compare(strings.ToLower(a.Borrower), strings.ToLower(b.Borrower))
		case models.IdxLoanedAt:
			r = a.Loaned.LoanedAt.Compare(b.LoanedAt)
		case models.IdxRating:
			r = cmp.Compare(a.Rating, b.Rating)
		case models.IdxCompletedAt:
			r = a.CompletedAt.Compare(b.CompletedAt)
		}
		if desending {
			return r * -1
		} else {
			return r
		}
	})
	return nil
}

// load load entries form repostory, sort them, and put them into table.
func (t *Table) load() error {

	items, err := t.repo.GetAllBooks()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	by, _ := t.SortBy.Get()
	order, _ := t.SortOrder.Get()
	err = sortBooks(items, by, order == "DESC")
	if err != nil {
		return err
	}
	for _, item := range items {
		err := t.table.AppendRow(
			item.ID,
			entryValues(&item),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) Size() (int, int) {
	return t.table.Size()
}

func (t *Table) Get(row, col int) string {
	cell := t.table.GetCell(row, col)
	value := cell.Value()
	if value == "" {
		value = stubValue
	}
	return value
}

func (t *Table) GetID(row, col int) (int64, error) {
	cell := t.table.GetCell(row, col)
	return cell.ID(), nil
}

// IsItemHidden check whether cell item is hidden.
func (t *Table) IsItemHidden(row, col int) bool {
	return t.table.GetCell(row, col).IsHidden()
}

// IsHeaderHidden check whether header at col is hidden.
func (t *Table) IsHeaderHidden(col int) bool {
	return t.table.GetHeader(col).IsHidden()

}

func (t *Table) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}


type TableControllersVM struct {
	SearchText    binding.String
	selector      *EntrySelect
	hiddenColumns []string
	bus           *bus.Bus
	EditIsOpen    binding.Bool
	editBook      *EditBookVM

	store repo.BookStore
}

func NewTableControllersVM(b *bus.Bus, r repo.BookRetriever, s repo.BookStore, g *UniqueGenres) *TableControllersVM {
	tc := &TableControllersVM{
		SearchText:    binding.NewString(),
		hiddenColumns: make([]string, 0),
		EditIsOpen:    binding.NewBool(),

		bus:           b,

		selector: newEntrySelect(r),
		store:    s,
	}
	tc.editBook = NewEditBookVM(b, s, g, tc.EditIsOpen)
	return tc
}

func msgNotSelected(b *bus.Bus) {
	b.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Nothing selected",
	})
}

func (tc *TableControllersVM) Delete() {
	if !tc.selector.HasSelected() {
		msgNotSelected(tc.bus)
		return
	}
	err := tc.store.DeleteBook(tc.selector.getID())
	if err != nil {
		log.Println(err)
		return
	}
	tc.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
}

func (tc *TableControllersVM) Edit() {
	if !tc.selector.HasSelected() {
		msgNotSelected(tc.bus)
		return
	}

	book, err := tc.selector.getBook()
	if err != nil {
		log.Println(err)
	}

	tc.editBook.reset()
	tc.editBook.Set(book)
	_ = tc.EditIsOpen.Set(true)
}

func (tc *TableControllersVM) Selector() *EntrySelect {
	return tc.selector
}

func (tc *TableControllersVM) GetEditBook() *EditBookVM {
	return tc.editBook
}

// entryHeaders lists the headers labels for book entry.
func entryHeaders() []string {
	return models.BookEntryFields()
}

// entryValues get the values from e in its in order of entryHeaders.
func entryValues(e *models.BookEntry) []string {
	
	headers := models.BookEntryFields()
	values  := make([]string, len(headers))

	for i, header := range headers {
		switch i {
		case models.IdxTitle:
			values[i] = e.Title
		case models.IdxAuthor:
			values[i] = e.Author
		case models.IdxGenre:
			values[i] = e.Genre
		case models.IdxRating:
			values[i] = formatRating(e.Rating)
		case models.IdxCompletedAt:
			values[i] = formatDate(&e.CompletedAt)
		case models.IdxBorrower:
			values[i] = e.Borrower
		case models.IdxLoanedAt:
			values[i] = formatDate(&e.LoanedAt)
		default:
			panic("unknown field:" + header)
		}
	}
	return values
}
