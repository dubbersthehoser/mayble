package table

import (
	"cmp"
	"log"
	"fmt"
	"slices"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
	"github.com/dubbersthehoser/mayble/internal/command"
)


type Point struct {
	Col int
	Row int
	ID  int64
}

type Table struct {
	Searching  *Searching

	Selected   *Selected
	Sorting    *Sorting
	Sheet      *Sheet
	Searchable *Searchable
	Settings   *Settings

	eb *event.EventBus
}


func NewTable(cfg *config.Config, cb *command.CommandBus, eb *event.EventBus) *Table {
	t := &Table{
		Sheet:      newSheet(eb, getShownHeader(cfg)),
		Searchable: newSearchable(getShownHeader(cfg), ColumnAll),
		Searching:  newSearching(cb),
		Sorting:    newSorting(cb, cfg.UI.TableSortBy, cfg.UI.TableAscending),
		Settings:   newSettings(eb, cfg),
		Selected:   newSelected(eb, cb),
	}
	return t
}

func SetupCommands(t *Table, eb *event.EventBus, cb *command.CommandBus) {
	cb.Register(CommandSnapshotSelect{}, func(v command.Command) error{

		e := v.(CommandSnapshotSelect)

		pp := snapshot.Current.Load()
		if pp.Version() != e.Version {
			return nil
		}

		if !e.Has {
			eb.Notify(EventSelected{
				Has: e.Has,
			})
			return nil
		}

		p, err := toSheetPoint(t.Sheet.Header(), e.Point, t.Sheet.IDToRow)
		if err != nil {
			return err
		}

		eb.Notify(EventSelected {
			Has: e.Has,
			Point: p,
		})
		return nil
	})
	cb.Register(CommandSearch, func(v command.Command) error {
		e := v.(CommandSearch)
	}
}


//
// Sheet
//

// Sheet a refrence view for table. 
// Methods should only be called by UI thread.
type Sheet struct {
	header       []string
	sorted       []int64
	idToRow      map[int64]int
	OnNewHeaders func()
	OnSorted     func()
}

func newSheet(eb *event.EventBus, header []string) *Sheet {
	s := &Sheet{
		header: header,
		sorted: make([]int64, 0),
		OnSorted: func() {},
		OnNewHeaders: func() {},
	}
	eb.Subscribe(EventColumnHidden{}, func(v event.Event) {
		e := v.(EventColumnHidden)
		header := make([]string, 0)
		for i, label := range models.BookEntryFields() {
			if !e.hidden[i] {
				header = append(header, label)
			}
		}
		s.header = header
		s.OnNewHeaders()
	})
	eb.Subscribe(EventSorted{}, func(v event.Event){
		ss := snapshot.Current.Load()
		e := v.(EventSorted)
		if ss.Version() == e.snapshot.Version() {
			s.sorted = e.ids
			s.OnSorted()
		}
	})
	return s
}

func (s *Sheet) RowToID(row int) (int64, error) {
	if len(s.sorted) <= row || row < 0 {
		return 0, fmt.Errorf("row_to_id %d: index out of range", row)
	}
	id := s.sorted[row]
	return id, nil
}

func (s *Sheet) IDToRow(id int64) (int, error) {
	row, ok := s.idToRow[id]
	if !ok {
		return 0, fmt.Errorf("id_to_row %d: id not found", id)
	}
	return row, nil
}

func (s *Sheet) Get(p Point) (string, error) {
	ss := snapshot.Current.Load()
	ssp, err := toSnapshotPoint(s.header, s.sorted, p, ss.IDToRow)
	if err != nil {
		return "", err
	}
	return ss.Get(ssp)
}

func (s *Sheet) Size() (rows, cols int) {
	ss := snapshot.Current.Load()
	rows, _  = ss.Size()
	cols = len(s.header)
	return
}

func (s *Sheet) Header() []string {
	return s.header
}

//
// Sorting
//

type Sorting struct {
	cb     *command.CommandBus
	Column    string
	Ascending bool
}

func newSorting(cb *command.CommandBus, column string, asc bool) *Sorting {
	s := &Sorting{
		cb: cb,
		Column: column,
		Ascending: asc,
	}
	return s
}

func (s *Sorting) Sort() {
	s.cb.Dispatch(CommandSort{
		asc: s.Ascending,
		column: s.Column,
	})
}

//
// Searchiable
//

const ColumnAll = "All"

type Searchable struct {
	headers  []string
	Selected string
	OnUpdate func()
}

func newSearchable(headers []string, by string) *Searchable {
	s := &Searchable{
		headers: headers,
		OnUpdate: func() {},
		Selected: by,
	}
	return s
}

func (s *Searchable) setSelectable(headers []string) {
	s.headers = headers
	s.OnUpdate()
}

func (s *Searchable) Options() []string {
	o := []string{
		ColumnAll,
	}
	o = append(o, s.headers...)
	return o
}

//
// Searching
//

type Searching struct {
	cb        *command.CommandBus
	debouncer *search.Debouncer
}

func newSearching(cb *command.CommandBus) *Searching {
	sr := &Searching{
		cb: cb,
	}
	return sr
}

func (s *Searching) Search(pattern string) {
	s.cb.Dispatch(CommandSearch{pattern: pattern})
}

//
// Selected
//

// Selected 
type Selected struct {
	cb        *command.CommandBus
	eb        *event.EventBus
	selected  Point
	has       bool
	OnChanged func()
}

func newSelected(eb *event.EventBus, cb *command.CommandBus) *Selected {
	es := &Selected{
		cb: cb,
		eb: eb,
		OnChanged: func(){},
		has: false,
	}
	es.eb.Subscribe(EventSelected{}, func(v event.Event) {
		e := v.(EventSelected)
		es.has = e.Has
		es.selected = e.Point
		es.OnChanged()
	})
	return es
}

func (es *Selected) Set(p Point, ok bool) {
	es.selected = p
	es.has = ok
	es.cb.Dispatch(CommandSelect{point: p, has: ok})
}

func (es *Selected) Get() (Point, bool) {
	return es.selected, es.has
}

//func (es *Selected) NextSearched() {
//	if es.searchedRow == -1 {
//		return
//	}
//	es.searchedRow += 1
//	if es.searchedRow >= len(es.searched) {
//		es.searchedRow = 0
//	}
//	es.selected = es.searched[es.searchedRow]
//	es.has = true
//	es.onChanged()
//}
//
//func (es *Selected) PrevSearched() {
//	if es.searchedRow == -1 {
//		return
//	}
//
//	es.searchedRow -= 1
//	if es.searchedRow < 0 {
//		es.searchedRow = len(es.searched)-1
//	}
//	es.selected = es.searched[es.searchedRow]
//	es.has = true
//	es.onChanged()
//}

//
// Settings
//

type Settings struct {
	eb       *event.EventBus
	cfg      *config.Config
}

func newSettings(eb *event.EventBus, cfg *config.Config) *Settings {
	cs := &Settings{
		cfg: cfg,
	}
	cs.eb.Subscribe(eventSorted{}, func(v event.Event) {
		e := v.(eventSorted)
		cfg.UI.TableSortBy = e.column
		cfg.UI.TableAscending = e.asc
	})
	return cs
}

func (ts *Settings) HeaderMinWidth() float32 {
	return ts.cfg.UI.TableMinWidth
}

func (ts *Settings) HeaderHeight() float32 {
	return ts.cfg.UI.TableHeaderHeight
}

func (ts *Settings) HeaderGetWidth(label string) float32 {
	idx := slices.Index(models.BookEntryFields(), label)
	return ts.cfg.UI.Headers[idx].Width
}

func (ts *Settings) HeaderSetWidth(label string, width float32) {
	idx := slices.Index(models.BookEntryFields(), label)
	h := ts.cfg.UI.Headers[idx]
	h.Width = width
	ts.cfg.UI.Headers[idx] = h
}

func (ts *Settings) IsLoanHidden() bool {
	return isLoanHidden(ts.cfg)
}

func (ts *Settings) IsReadHidden() bool {
	return isReadHidden(ts.cfg)
}

func (ts *Settings) IsIDHidden() bool {
	return isIDHidden(ts.cfg)
}

func (ts *Settings) SetIDHidden(t bool) {
	header := ts.cfg.UI.Headers[models.IdxID]
	header.IsHidden = t
	ts.cfg.UI.Headers[models.IdxID] = header
	ts.notifyHidden()
}

func (ts *Settings) ToggleHiddenID() bool {
	header := ts.cfg.UI.Headers[models.IdxID]
	t := !header.IsHidden
	ts.SetIDHidden(t)
	return t
}

func (ts *Settings) SetLoanHidden(t bool) {
	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]

	loaned.IsHidden = t
	borrower.IsHidden = t

	ts.cfg.UI.Headers[models.IdxLoanedAt] = loaned
	ts.cfg.UI.Headers[models.IdxBorrower] = borrower
	ts.notifyHidden()
}

func (ts *Settings) ToggleHiddenLoan() bool {
	loaned := ts.cfg.UI.Headers[models.IdxLoanedAt]
	borrower := ts.cfg.UI.Headers[models.IdxBorrower]
	t := !(loaned.IsHidden || borrower.IsHidden)
	ts.SetLoanHidden(t)
	return t
}

func (ts *Settings) SetReadHidden(t bool) {
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]

	rating.IsHidden = t
	completed.IsHidden = t

	ts.cfg.UI.Headers[models.IdxRating] = rating
	ts.cfg.UI.Headers[models.IdxCompletedAt] = completed
	ts.notifyHidden()
}

func (ts *Settings) ToggleHiddenRead() bool {
	rating := ts.cfg.UI.Headers[models.IdxRating]
	completed := ts.cfg.UI.Headers[models.IdxCompletedAt]
	t := !(completed.IsHidden || rating.IsHidden)
	ts.SetReadHidden(t)
	return t
}

func (ts *Settings) SetWidth(label string, width float32) {
	idx := slices.Index(models.BookEntryFields(), label)
	if idx == -1 {
		log.Printf("Error: invalid header label '%s'", label)
		return
	}
	h, ok := ts.cfg.UI.Headers[idx]
	if !ok {
		log.Printf("Warning: column not found '%s'", label)
		return
	}
	h.Width = width
	ts.cfg.UI.Headers[idx] = h
}

func (ts *Settings) GetWidth(label string) float32 {
	idx := slices.Index(models.BookEntryFields(), label)
	if idx == -1 {
		log.Printf("Error: invalid header label '%s'", label)
		return ts.cfg.UI.TableMinWidth
	}
	h := ts.cfg.UI.Headers[idx]
	width := h.Width
	return width
}

func (ts *Settings) notifyHidden() {
	hidden := make([]bool, len(ts.cfg.UI.Headers))
	for idx := range models.BookEntryFields() {
		hidden[idx] = ts.cfg.UI.Headers[idx].IsHidden
	}
	ts.eb.Notify(eventHiddenColumn{
		hidden: hidden,
	})
}

//
// Helper Functions
//

//func snapshotSearch(result chan <- searchResult, ss *snapshot.Snapshot, by string, pattern string) {
//	var trv search.Traverser
//	if by == ColumnAll {
//		trv = newTableTraverse(ss)
//	} else {
//		idx := slices.Index(models.BookEntryFields(), by)
//		if idx == -1 {
//			log.Printf("search %s: invalid column label", by)
//			return
//		}
//		trv = newColumnTraverse(ss, idx)
//	}
//	srch := (&search.Searcher{}).Set(trv, pattern)
//	points, score := searchSearcher(srch)
//	result <- searchResult{
//		Points: points,
//		Version: ss.Version(),
//		Score: score,
//	}
//
//}

func searchSearcher(srch *search.Searcher) ([]Point, []int) {
	type result struct {
		row,
		col,
		score int
		id    int64
	}
	results := make([]result, 0)
	for srch.Next() {
		point := srch.Point()
		score := srch.Score()
		if score == -1 {
			continue
		}
		r := result{
			row:   point.Row,
			col:   point.Col,
			score: srch.Score(),
		}
		results = append(results, r)
	}
	if len(results) == 0 {
		return []Point{}, []int{}
	}
	slices.SortFunc(results, func(a, b result) int {
		r := cmp.Compare(a.score, b.score)
		if r == 0 {
			return cmp.Compare(a.row, b.row)
		}
		return r * -1
	})
	points := make([]Point, len(results))
	score := make([]int, len(results))
	for i, r := range results {
		p := Point{
			Row: r.row,
			Col: r.col,
			ID: r.id,
		}
		points[i] = p
		score[i] = r.score
	}
	return points, score
}

func isLoanHidden(cfg *config.Config) bool {
	loaned := cfg.UI.Headers[models.IdxLoanedAt]
	borrower := cfg.UI.Headers[models.IdxBorrower]
	return loaned.IsHidden && borrower.IsHidden

}

func isIDHidden(cfg *config.Config) bool {
	header := cfg.UI.Headers[models.IdxID]
	return header.IsHidden
}

func isReadHidden(cfg *config.Config) bool {
	rating := cfg.UI.Headers[models.IdxRating]
	completed := cfg.UI.Headers[models.IdxCompletedAt]
	return rating.IsHidden && completed.IsHidden
}


func getShownHeader(cfg *config.Config) []string {	
	set := make([]string, 0)
	for _, h := range cfg.UI.Headers {
		if !h.IsHidden {
			set = append(set, h.Name)
		}
	}
	return set
}

func sortIDs(ids []int64, ss *snapshot.Snapshot, column string, asc bool) {
	colIdx := slices.Index(models.BookEntryFields(), column)
	comp, err := app.CompareBookEntry(colIdx, asc)
	if err != nil {
		log.Println("sorting:", err)
		return
	}
	slices.SortFunc(ids, func(a, b int64) int {
		bookA, _ := ss.GetBookEntryByID(a)
		bookB, _ := ss.GetBookEntryByID(b)
		return comp(*bookA, *bookB)
	})
}

func isValidVersion(curr, result *snapshot.Snapshot) bool {
	return curr != nil && curr.Version() != result.Version()
}

func toSnapshotPoint(
	header []string, 
	sorted []int64, 
	p Point, 
	getRowByID func(int64) (int, error),
) (snapshot.Point, error) {

	col, err := toSnapshotColumn(header, p.Col)
	if err != nil {
		return snapshot.Point{}, err
	}
	row, err := toSnapshotRow(sorted, p.Row, getRowByID)
	if err !=nil {
		return snapshot.Point{}, err
	}
	return snapshot.Point{
		Row: row,
		Col: col,
	}, err
}

func toSnapshotColumn(header []string, column int) (int, error) {
	if column >= len(header) || column < 0 {
		return 0, fmt.Errorf("to_snapshot_column %d: index out of range", column)
	}
	return slices.Index(models.BookEntryFields(), header[column]), nil
}

func toSnapshotRow(sorted []int64, row int, getRowByID func(int64) (int, error)) (int, error) {
	if row >= len(sorted) || row < 0 {
		return 0, fmt.Errorf("to_snapshot_row %d: index out of range", row)
	}
	id := sorted[row]
	return getRowByID(id)
}

func toSheetPoint(
	header []string,
	p snapshot.Point,
	getRowByID func(int64) (int, error),
) (Point, error) {

	col, err := toSheetColumn(header, p.Col)
	if err != nil {
		return Point{}, err
	}

	row, err := getRowByID(p.ID)
	if err != nil {
		return Point{}, err
	}

	return Point{
		Row: row,
		Col: col,
	}, err
}

func toSheetColumn(header []string, column int) (int, error) {
	if column >= len(models.BookEntryFields()) || column < 0 {
		return 0, fmt.Errorf("to_sheet_column %s: index out of range", column)
	}
	label := models.BookEntryFields()[column]
	col := slices.Index(header, label)
	if col == -1 {
		return 0, fmt.Errorf("to_sheet_column %s: label not found", label)
	}
	return col, nil
}


