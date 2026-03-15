package viewmodel

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strings"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/table"
)

const (
	notifySelect bool   = true
	stubValue    string = "N/A"
)

type TableVM struct {
	repo   repo.BookRetriever
	config *config.Config
	bus    *bus.Bus

	SortBy    binding.String
	SortOrder binding.String

	selector *EntrySelect

	Search struct {
		Text   binding.String
		Header binding.String
	}

	table *table.Table

	l *listener
}

func NewTableVM(b *bus.Bus, app *appService) *TableVM {
	t := &TableVM{
		table:  table.NewTable("Main", entryHeaders()),
		repo:   app.bookRetriever,
		config: app.cfg,
		bus:    b,

		SortBy:    binding.NewString(),
		SortOrder: binding.NewString(),

		Search: struct {
			Text   binding.String
			Header binding.String
		}{
			Text:   binding.NewString(),
			Header: binding.NewString(),
		},

		selector: newEntrySelect(app.bookRetriever),

		l: &listener{},
	}

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

	if t.config != nil {
		t.table.SetHidden(t.config.UI.Table.ColumnsHidden)
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

func (t *TableVM) search() {
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

func (t *TableVM) SetSelector(es *EntrySelect) *TableVM {
	t.selector = es
	return t
}

// SearchOptions a list of searchable options.
func (t *TableVM) SearchOptions() []string {
	return []string{
		"All",
		"Title",
		"Author",
		"Genre",
		"Borrower",
	}
}

// Selector returns the table's selector.
func (t *TableVM) Selector() *EntrySelect {
	return t.selector
}

// Sort table using sort bindings.
func (t *TableVM) Sort() error {
	t.selector.unselect(notifySelect)
	err := t.reload()
	if err != nil {
		return err
	}
	t.l.notify()
	return nil
}

// The smallest width that a column can be.
const MinColWidth float32 = 100.0

// StoreColumnWidth to the config file if it exists, else nop.
// When width is smaller then MinColWidth, MinColWidth will be used.
func (t *TableVM) StoreColumnWidth(col int, width float32) {
	if t.config == nil {
		return
	}
	if width < MinColWidth {
		width = MinColWidth
	}
	table := t.config.GetUITable()
	cell := t.table.GetHeaderCell(col)
	label := cell.Header()
	table.SetColumnWidth(label, width)
}

// GetColumnWidth from the config file if it exsits, else returns defualt MinColWidth.
func (t *TableVM) GetColumnWidth(col int) float32 {
	if t.config == nil {
		return MinColWidth
	}
	cell := t.table.GetHeaderCell(col)
	label := cell.Header()

	var width float32 = 0.0
	if t.config != nil {
		table := t.config.GetUITable()
		width = table.GetColumnWidth(label)
	}
	if width < MinColWidth {
		width = MinColWidth
	}
	return width
}

func (t *TableVM) SetHidden(options []string) {
	hide := hiddenOptionsToHeaders(options)
	if t.config != nil {
		table := t.config.GetUITable()
		table.ColumnsHidden = hide
	}
	t.table.SetHidden(hide)
	t.l.notify()
}

func (t *TableVM) Hidden() []string {
	headers := t.table.HiddenHeaders()
	return hiddenHeadersToOptions(headers)
}

// hiddenHeadersToOptions returns hidden options from headers.
func hiddenHeadersToOptions(headers []string) []string {
	options := slices.Clone(headers)
	options = slices.DeleteFunc(options, func(s string) bool {
		return s == "Rating" || s == "Borrower"
	})
	return options
}

// hiddenOptionsToHeaders returns a slice of headers from a slice of hidden header options.
func hiddenOptionsToHeaders(options []string) []string {
	hide := make([]string, 0)
	for _, o := range options {
		switch o {
		case "Loaned":
			hide = append(hide, "Loaned", "Borrower")
		case "Read":
			hide = append(hide, "Read", "Rating")
		default:
			hide = append(hide, o)
		}
	}
	return hide
}

// HiddenOptions returns the list of options for hiding columns.
func (t *TableVM) HiddenOptions() []string {
	headers := repo.BookEntryFields()
	return hiddenHeadersToOptions(headers)
}

func (t *TableVM) Headers() []string {
	return t.table.Headers()
}

func (t *TableVM) Select(row, col int) {
	cell := t.table.GetCell(row, col)
	t.selector.selectID(cell.ID(), !notifySelect)
	t.selector.selectCell(row, col, notifySelect)
}

func (t *TableVM) Unselect(row, col int) {
	t.selector.unselect(!notifySelect)
}

// reload clear table, then call load.
func (t *TableVM) reload() error {
	t.table.ClearValues()
	return t.load()
}

// sortbooks sort slice of books.
func sortBooks(books []repo.BookEntry, header string, desending bool) error {

	index := slices.Index(entryHeaders(), header)

	if index == -1 {
		return fmt.Errorf("sort_books: invalid header '%s'", header)
	}

	slices.SortFunc(books, func(a, b repo.BookEntry) int {
		r := -1
		switch index {
		case repo.IdxTitle:
			r = cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
		case repo.IdxAuthor:
			r = cmp.Compare(strings.ToLower(a.Author), strings.ToLower(b.Author))
		case repo.IdxGenre:
			r = cmp.Compare(strings.ToLower(a.Genre), strings.ToLower(b.Genre))
		case repo.IdxBorrower:
			r = cmp.Compare(strings.ToLower(a.Borrower), strings.ToLower(b.Borrower))
		case repo.IdxLoaned:
			r = a.Loaned.Compare(b.Loaned)
		case repo.IdxRating:
			r = cmp.Compare(a.Rating, b.Rating)
		case repo.IdxRead:
			r = a.Read.Compare(b.Read)
		}
		if desending {
			return r * -1
		} else {
			return r
		}
	})
	return nil
}

// load load entries form repostory sort then, and put them into table.
func (t *TableVM) load() error {

	items, err := t.repo.GetAllBooks(repo.Book)
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

func (t *TableVM) Size() (int, int) {
	return t.table.Size()
}

func (t *TableVM) Get(row, col int) string {
	cell := t.table.GetCell(row, col)
	value := cell.Value()
	if value == "" {
		value = stubValue
	}
	return value
}

func (t *TableVM) GetID(row, col int) (int64, error) {
	cell := t.table.GetCell(row, col)
	return cell.ID(), nil
}

// IsItemHidden check whether cell item is hidden.
func (t *TableVM) IsItemHidden(row, col int) bool {
	cell := t.table.GetCell(row, col)
	return t.table.IsHidden(cell)
}

// IsHeaderHidden check whether header at col is hidden.
func (t *TableVM) IsHeaderHidden(col int) bool {
	cell := t.table.GetHeaderCell(col)
	return t.table.IsHidden(cell)

}

func (t *TableVM) AddListener(l binding.DataListener) {
	t.l.AddListener(l)
}

type TableControllersVM struct {
	SearchText    binding.String
	selector      *EntrySelect
	hiddenColumns []string
	app           *appService
	bus           *bus.Bus
	EditIsOpen    binding.Bool
	editBook      *EditBookVM
	table         *TableVM
}

func NewTableControllersVM(b *bus.Bus, app *appService) *TableControllersVM {
	tc := &TableControllersVM{
		SearchText:    binding.NewString(),
		hiddenColumns: make([]string, 0),
		selector:      newEntrySelect(app.bookRetriever),
		EditIsOpen:    binding.NewBool(),
		app:           app,
		bus:           b,
	}
	tc.editBook = NewEditBookVM(b, app, tc.EditIsOpen)
	return tc
}

func (tc *TableControllersVM) Delete() {
	if !tc.selector.HasSelected() {
		tc.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "Nothing selected",
		})
		return
	}
	book, err := tc.selector.getBook()
	if err != nil {
		log.Println(err)
		return
	}
	err = tc.app.bookDeletor.DeleteBook(book)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(book)
	tc.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
}
func (tc *TableControllersVM) Edit() {
	if !tc.selector.HasSelected() {
		tc.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "Nothing selected",
		})
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
	return repo.BookEntryFields()
}

// entryValues get the values from e in its in order of entryHeaders.
func entryValues(e *repo.BookEntry) []string {
	return []string{
		e.Title,
		e.Author,
		e.Genre,
		formatRating(e.Rating),
		formatDate(&e.Read),
		e.Borrower,
		formatDate(&e.Loaned),
	}
}
