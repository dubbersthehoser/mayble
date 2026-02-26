package viewmodel


import (
	"slices"
	"cmp"
	"fmt"
	"log"

	"fyne.io/fyne/v2/data/binding"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/table"
)



type TableVM struct {
	repo     repo.BookRetriever
	config   *config.Config

	SortBy     binding.String
	SortOrder  binding.String

	selector   *EntrySelect

	Search   *TableSearch
	
	table    *table.Table
	
	l        *listener
}


func NewTableVM(vms *vmService) *TableVM {
	t := &TableVM{
		table:   table.NewTable("Main", entryHeaders()),
		repo:    vms.app.bookRetriever,
		config:  vms.app.cfg,

		SortBy:     binding.NewString(),
		SortOrder:  binding.NewString(),

		selector: newEntrySelect(vms.app.bookRetriever),

		l: &listener{},
	}

	t.Search = NewTableSearch(t.table)

	t.Search.Text.AddListener(binding.NewDataListener(func() {
		t.selector.unselect(true)
		search, _ := t.Search.Text.Get()
		if search == "" {
			return
		}
		result := t.Search.search(search)
		if len(result) == 0 {
			return
		}
		r := result[0]
		t.selector.selectID(r.id, false)
		t.selector.selectCell(r.row, r.col, true)
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

	vms.bus.Subscribe(bus.Handler{
		Name: msgDataChanged,
		Handler: func(e *bus.Event) {
			t.table.ClearValues()
			t.selector.unselect(true)
			err := t.load()
			if err != nil {
				log.Println(err)
				return
			}
			t.l.notify()
		},
	})

	return t
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
func (t *TableVM) Sort() {
	err := t.table.ClearValues()
	if err != nil {
		log.Println(err)
		return
	}
	t.selector.unselect(true)
	err = t.load()
	if err != nil {
		log.Println(err)
		return
	}
	t.l.notify()
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

func hiddenHeadersToOptions(hide []string) []string {
	options := make([]string, 0)
	hasLoaned := slices.Contains(hide, "Borrower") && slices.Contains(hide, "Loaned")
	hasRead := slices.Contains(hide, "Rating") && slices.Contains(hide, "Read")
	for _, h := range hide {
		switch h {
		case "Borrower", "Loaned", "Rating", "Read":
			continue
		default:
			options = append(options, h)
		}
	}
	if hasLoaned {
		options = append(options, "Loaned")
	}
	if hasRead {
		options = append(options, "Read")
	}
	return options
}

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


func (t *TableVM) HideOptions() []string {
	return []string{
		"Title",
		"Author",
		"Genre",
		"Read",
		"Loaned",
	}
}

func (t *TableVM) Headers() []string {
	return t.table.Headers()
}


func (t *TableVM) Select(row, col int) {
	cell := t.table.GetCell(row, col)
	t.selector.selectID(cell.ID(), false)
	t.selector.selectCell(row, col, true)
}

func (t *TableVM) Unselect(row, col int) {
	t.selector.unselect(false)
}


// load entries from repository sort them into table.
func (t *TableVM) load() error {

	items, err := t.repo.GetAllBooks(repo.BookLoaned | repo.BookRead)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	// This should be part of application,
	// but whatever...
	by, _ := t.SortBy.Get()
	order, _ := t.SortOrder.Get()

	index := slices.Index(entryHeaders(), by)

	slices.SortFunc(items, func(a, b repo.BookEntry) int {
		r := -1
		switch index {
		case 0:
			r = cmp.Compare(a.Title, b.Title)
		case 1:
			r = cmp.Compare(a.Author, b.Author)
		case 2:
			r = cmp.Compare(a.Genre, b.Genre)
		case 3:
			r = cmp.Compare(a.Borrower, b.Borrower)
		case 4:
			r = a.Loaned.Compare(b.Loaned)
		case 5:
			r = cmp.Compare(a.Rating, b.Rating)
		case 6:
			r = a.Read.Compare(b.Read)
		default:
			log.Println("load: sort field not found", index, by)
		}
		if order == "DESC" {
			return r * -1
		} else {
			return r
		}
	})

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
		value = "N/A"
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
	vms           *vmService
	EditIsOpen    binding.Bool
	editbook      *EditBookVM            
	table         *TableVM
}

func NewTableControllersVM(vms *vmService) *TableControllersVM {
	tc := &TableControllersVM{
		SearchText: binding.NewString(),
		hiddenColumns: make([]string, 0),
		selector: newEntrySelect(vms.app.bookRetriever),
		EditIsOpen: binding.NewBool(),
		vms: vms,
	}
	tc.editbook = NewEditBookVM(vms, tc.EditIsOpen)
	return tc
}

func (tc *TableControllersVM) Delete() {
	if tc.selector.HasSelected() {
		book, err := tc.selector.getBook()
		if err != nil {
			log.Println(err)
			return
		}
		err = tc.vms.app.bookDeletor.DeleteBook(book)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(book)
		tc.vms.bus.Notify(bus.Event{
			Name: msgDataChanged,
		})
	} else {
		tc.vms.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "Nothing selected",
		})
	}
}
func (tc *TableControllersVM) Edit() {
	if !tc.selector.HasSelected() {
		tc.vms.bus.Notify(bus.Event{
			Name: msgUserInfo,
			Data: "Nothing selected",
		})
		return
	}
	book, err := tc.selector.getBook()
	if err != nil {
		log.Println("edit:", err)
	}
	tc.editbook.reset()
	tc.editbook.Set(book)
	fmt.Println(book)
	_ = tc.EditIsOpen.Set(true)
}

func (tc *TableControllersVM) Selector() *EntrySelect {
	return tc.selector
}

func (tc *TableControllersVM) GetEditBook() *EditBookVM {
	return tc.editbook
}


// entryHeaders lists the headers labels for book entry.
func entryHeaders() []string {
	return []string{
		"Title",
		"Author",
		"Genre",

		"Rating",
		"Read",

		"Borrower",
		"Loaned",
	}
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
