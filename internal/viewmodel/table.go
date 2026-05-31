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
	"github.com/dubbersthehoser/mayble/internal/config"
)

const (
	notifySelect bool   = true
	stubValue    string = "N/A"
)

const (
	tableEntrySelected   string = "table.selected"
	tableEntryUnselected string = "table.unselected"
	tableEditOpen        string = "table.edit.open"
	tableEditClose       string = "table.edit.close"
)

// The smallest width a column can be.
const MinColWidth float32 = 100.0

type TableConfig struct {
	cfg *config.Config
}

func (tc *TableConfig) GetHiddenColumns() []string {
	columns := make([]string, 0)
	for label, headers := range tc.cfg.UI.Headers {
		if headers.IsHidden {
			columns = append(columns, label)
		}
	}
	return columns
}

func (tc *TableConfig) SetHiddenColumns(labels []string) {
	for _, label := range labels {
		header, ok := tc.cfg.UI.Headers[label]
		if !ok {
			continue
		}
		header.IsHidden = true
		tc.cfg.UI.Headers[label] = header
	}
}

func (tc *TableConfig) SetSortBy(label string) {
	tc.cfg.UI.TableSortBy = label
}
func (tc *TableConfig) GetSortBy() string {
	by := tc.cfg.UI.TableSortBy
	if by == "" {
		by = models.BookEntryFields()[models.IdxTitle]
		tc.SetSortBy(by)
	}
	return by
}

func (tc *TableConfig) GetAscending() bool {
	return tc.cfg.UI.TableAscending
}

func (tc *TableConfig) SetAscending(t bool) {
	tc.cfg.UI.TableAscending = t
}

func (tc *TableConfig) SetColumnWidth(label string, size float32) {
	if size < MinColWidth {
		size = MinColWidth
	}
	header, ok := tc.cfg.UI.Headers[label]
	if !ok {
		return
	}
	header.Width = size
	tc.cfg.UI.Headers[label] = header
}

func (tc *TableConfig) GetColumnWidth(label string) float32 {
	header, ok := tc.cfg.UI.Headers[label]
	if !ok {
		return MinColWidth
	}
	size := header.Width
	if size < MinColWidth {
		size = MinColWidth
	}
	return size
}

type Table struct {
	store     repo.BookStore
	retriever repo.BookRetriever

	bus    *bus.Bus
	cfg    *TableConfig
	genres *UniqueGenres

	table *table.Table

	l *listener
}

func NewTable(b *bus.Bus, cfg *TableConfig, source sourceSubject, s repo.BookStore, r repo.BookRetriever, ug *UniqueGenres) *Table {
	t := &Table{
		table: table.NewTable("Main", entryHeaders()),
		store:  s,
		retriever: r,
		cfg:   cfg,
		bus:   b,
		genres: ug,

		l: &listener{},
	}
	source.AddListener(func(){
		err := t.reload()
		if err != nil {
			log.Println(err)
			return
		}
		t.l.notify()
	})
	return t
}

// StoreColumnWidth to the config file if it exists else it will be an nop.
// When width is smaller then MinColWidth, MinColWidth will be used.
func (t *Table) StoreColumnWidth(col int, width float32) {
	header := t.table.GetHeader(col)
	label := header.Name()
	t.cfg.SetColumnWidth(label, width)
}

// GetColumnWidth from the config file if it exsits, else returns defualt MinColWidth.
//func (t *Table) GetColumnWidth(col int) float32 {
//	label := t.table.GetHeader(col).Name()
//	width := t.cfg.GetColumnWidth(label)
//	if width < MinColWidth {
//		width = MinColWidth
//	}
//	return width
//}

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

func (t *Table) IsHidden(row, col int) bool {
	cell := t.table.GetCell(row, col)
	return cell.IsHidden()
}

func (t *Table) GetID(row, col int) (int64, error) {
	cell := t.table.GetCell(row, col)
	return cell.ID(), nil
}

func (t *Table) Sync() {
	t.reload()
}

func (t *Table) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}

// reload clear table, then call load.
func (t *Table) reload() error {
	t.table.ClearValues()
	return t.load()
}

// load load entries form repostory, sort them, and put them into table.
func (t *Table) load() error {

	items, err := t.retriever.GetAllBooks()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	by := t.cfg.GetSortBy()
	ascending := t.cfg.GetAscending()
	err = sortBooks(items, by, ascending)
	if err != nil {
		return err
	}
	return nil
}


type TableHeaders struct {
	table   *Table

	// hiding columns
	options []string
	aliased map[string][]string

	Labels      map[string]binding.String
	labelSuffix map[string]string
}

func NewTableHeaders(table *Table) *TableHeaders {
	columns := models.BookEntryFields()
	h := TableHeaders{
		table: table,
		options: []string{"Title", "Author", "Genre", "Read", "Loaned"},
		aliased: map[string][]string{
			"Title": {columns[models.IdxTitle]},
			"Author": {columns[models.IdxAuthor]},
			"Genre": {columns[models.IdxGenre]},
			"Read": {
				columns[models.IdxRating],
				columns[models.IdxCompletedAt],
			},
			"Loaned": {
				columns[models.IdxLoanedAt],
				columns[models.IdxBorrower],
			},
		},

		Labels: make(map[string]binding.String),

		labelSuffix: map[string]string{
			"normal": "- ",
			"asc": "↑ ",
			"desc": "↓ ",
		},
	}

	for _, label := range h.Headers() {
		h.Labels[label] = binding.NewString()
		_ = h.Labels[label].Set(h.labelSuffix["normal"] + label)
	}

	return &h
}

func (hh *TableHeaders) Headers() []string {
	return hh.table.table.Headers()
}

func (hh *TableHeaders) HideOptions() []string {
	return hh.options
}

func (hh *TableHeaders) SetHidden(options []string) {
	columns := make([]string, 0)
	for _, o := range options {
		cols, ok := hh.aliased[o]
		if !ok {
			log.Printf("WARNING: invalid option for hidden column '%s'", o)
			continue
		}
		columns = append(columns, cols...)
	}

	hh.table.cfg.SetHiddenColumns(columns)
	hh.table.table.SetHidden(columns)
	hh.table.l.notify()
}

func (hh *TableHeaders) GetHidden() []string {
	headers := hh.table.cfg.GetHiddenColumns()
	options := make([]string, 0)
	for option, aliased := range hh.aliased {
		for _, column := range aliased {
			if slices.Contains(headers, column) {
				options = append(options, option)
			}
		}
	}
	return options
}

func (hh *TableHeaders) IsHidden(col int) bool {
	header := hh.table.table.GetHeader(col)
	return header.IsHidden()
}

func (th *TableHeaders) Sort(label string) {
	by := th.table.cfg.GetSortBy()
	ascending := th.table.cfg.GetAscending()
	if by == label {
		ascending = !ascending
	} else {
		ascending = false
	}

	for _, label := range th.Headers() {
		th.Labels[label].Set(th.labelSuffix["normal"] + label)
	}



	if ascending {
		th.Labels[label].Set(th.labelSuffix["asc"] + label)
	} else {
		th.Labels[label].Set(th.labelSuffix["desc"] + label)
	}

	th.table.cfg.SetAscending(ascending)
	th.table.cfg.SetSortBy(label)


	books, err := th.table.retriever.GetAllBooks()
	if err != nil {
		log.Println(err)
		return
	}

	err = sortBooks(books, by, ascending)
	if err != nil {
		log.Println(err)
		return
	}

	th.table.table.ClearValues()

	err = loadBooksToTable(books, th.table.table)
	if err != nil {
		log.Println(err)
		return
	}
}

func (th *TableHeaders) SetWidthWithColumn(col int, width float32) {
	label := th.table.table.GetHeader(col).Name()
	th.table.cfg.SetColumnWidth(label, width)
}

func (th *TableHeaders) GetWidthWithLabel(label string) float32 {
	return th.table.cfg.GetColumnWidth(label)
}


type TableSelect struct {
	table       *Table
	cell        *table.Cell
	hasSelected bool

	searchOptions []string
	searchBy    string

	l *listener
}

func NewTableSelect(t *Table) *TableSelect {
	ts := &TableSelect{
		table: t,
		searchBy: "All",
		searchOptions: []string{
			"All",
			"Title",
			"Author",
			"Genre",
			"Borrower",
		},
		l: &listener{},
	}
	return ts
}

func (ts *TableSelect) Select(row, col int) {
	ts.hasSelected = true
	ts.l.notify()
	ts.cell = ts.table.table.GetCell(row, col)
	ts.table.bus.Notify(bus.Event{
		Name: tableEntrySelected,
		Data: ts.cell.ID(),
	})
}

func (ts *TableSelect) Unselect() {
	ts.hasSelected = false
	ts.table.bus.Notify(bus.Event{
		Name: tableEntryUnselected,
		Data: -1,
	})
}

func (ts *TableSelect) Selected() (row, col int) {
	return ts.cell.Point()
}

func (ts *TableSelect) HasSelected() bool {
	return ts.hasSelected
}

func (ts *TableSelect) Search(s string) {
	results := table.Search(ts.table.table, s, ts.searchBy)
	if len(results) == 0 {
		return
	}
	result := results[0]
	ts.hasSelected = true
	ts.cell = ts.table.table.GetCell(result.Row, result.Col)
	ts.l.notify()
	ts.table.bus.Notify(bus.Event{
		Name: tableEntrySelected,
		Data: ts.cell.ID(), 
	})
}

func (ts *TableSelect) SetSearchBy(option string) {
	ts.searchBy = option
}

func (ts *TableSelect) GetSearchBy() string {
	return ts.searchBy
}


func (ts *TableSelect) SearchOptions() []string {
	return ts.searchOptions
}

func (t *TableSelect) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}


type TableEdit struct {
	BookForm
	table   *Table
	IsOpen  binding.Bool
	selected *models.BookEntry
}

func NewTableEdit(t *Table) *TableEdit {
	ed := &TableEdit{
		BookForm: *NewBookForm(),
		IsOpen: binding.NewBool(),
		table: t,
	}

	// handle selected book events
	ed.table.bus.Subscribe(bus.Handler{
		Name: tableEntrySelected,
		Handler: func(e *bus.Event) {
			id := e.Data.(int64)
			book, err := ed.table.retriever.GetBookByID(id)
			if err != nil {
				log.Println(err)
				return
			}
			ed.selected = &book
		},
	})

	// handle un-selected events
	ed.table.bus.Subscribe(bus.Handler{
		Name: tableEntryUnselected,
		Handler: func(e *bus.Event) {
			ed.selected = nil
		},
	})

	return ed
}

func (ed *TableEdit) Genres() *UniqueGenres {
	return ed.table.genres
}

func (ed *TableEdit) Delete() {
	if ed.selected == nil {
		log.Println("ERROR: table edit selected is nil")
		return
	}
	ed.table.store.DeleteBook(ed.selected.ID)
	ed.table.reload()
}

func (ed *TableEdit) Open() {
	if ed.selected == nil {
		log.Println("ERROR: table edit selected is nil")
		return
	}
	ed.Set(ed.selected)
	_ = ed.IsOpen.Set(true)
}

func (ed *TableEdit) Submit() {
	err := ed.BookForm.validate()
	if err != nil {
		ed.table.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err.Error(),
		})
		return
	}

	book := ed.BookForm.ToBookEntry()

	ed.table.store.UpdateBook(book)
	ed.table.reload()
	ed.Close()
}

func (ed *TableEdit) Close() {
	ed.BookForm.reset()
	ed.IsOpen.Set(false)
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

// sortbooks sort slice of book entries.
func sortBooks(books []models.BookEntry, header string, ascending bool) error {

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
		if !ascending {
			return r * -1
		} else {
			return r
		}
	})
	return nil
}

func msgNotSelected(b *bus.Bus) {
	b.Notify(bus.Event{
		Name: msgUserInfo,
		Data: "Nothing selected",
	})
}

func loadBooksToTable(books []models.BookEntry, t *table.Table) error {
	for _, book := range books {
		err := t.AppendRow(
			book.ID,
			entryValues(&book),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
