package table

import (
	"cmp"
	"log"
	"fmt"
	"slices"
	"time"
	"context"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/worker"
)

const ColumnAll = "All"

type Point struct {
	Col int
	Row int
	ID  int64
}

type Table struct {
	worker     *worker.Worker
	Searching  *Searching
	Selected   *Selected
	Sorting    *Sorting
	Sheet      *Sheet
	Searchable      *Searchable
	SearchSelection *SearchSelection
	Settings   *Settings

	eb *event.EventBus
}


func NewTable(cfg *config.Config, w *worker.Worker, cb *command.CommandBus, eb *event.EventBus) *Table {
	t := &Table{
		worker:     w,
		Sheet:      newSheet(eb, getShownHeader(cfg)),
		Searchable: newSearchable(getShownHeader(cfg)),
		Searching:  newSearching(ColumnAll, cb),
		Sorting:    newSorting(cb, cfg.UI.TableSortBy, cfg.UI.TableAscending),
		Settings:   newSettings(eb, cfg),
		Selected:   newSelected(cb),
		SearchSelection: newSearchSelection(eb, cb),
	}
	SetupCommands(t, eb, cb)
	t.Searchable.onChangedSearchBy = t.Searching.setSearchColumn
	return t
}

func (t *Table) HandleWorkerEvent(ev worker.Event) {
	v, ok := ev.Data.(event.Event)
	if !ok {
		log.Println("Error: invalid worker event data")
		return
	}
	
	switch v.(type) {
	case EventSnapshotSorted, EventSnapshotSearched:
		t.eb.Notify(v)
	default:
		log.Println("Error: invalid worker event data")
	}
}

func SetupCommands(t *Table, eb *event.EventBus, cb *command.CommandBus) {

	// CommandSnapshotSelect
	cb.Register(CommandSnapshotSelect{}, func(v command.Command) error{

		e := v.(CommandSnapshotSelect)

		pp := snapshot.Current.Load()
		if pp.Version() != e.Version {
			log.Printf("Warning: de-synced versions: %d != %d", pp.Version(), e.Version)
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
			log.Println("Error:", err)
			return nil
		}

		eb.Notify(EventSelected {
			Has: e.Has,
			Point: p,
		})
		return nil
	})

	// CommandSheetSelect
	cb.Register(CommandSheetSelect{}, func(v command.Command) error {
		e := v.(CommandSheetSelect)

		if !e.Has {
			eb.Notify(EventSelected{
				Has: e.Has,
			})
			return nil
		}

		eb.Notify(EventSelected{
			Has:   e.Has,
			Point: e.Point,
		})

		return nil
	})

	// CommandSearch
	cb.Register(CommandSearch{}, func(v command.Command) error {
		e := v.(CommandSearch)
		pattern := e.Pattern
		column := e.Column
		t.worker.Jobs <- NewJobSearchSnapshot(t.worker, pattern, column) 
		return nil
	})

	// CommandSort
	cb.Register(CommandSort{}, func(v command.Command) error {
		e := v.(CommandSort)
		column := e.Column
		asc := e.Asc
		t.worker.Jobs <- NewJobSortSnapshot(t.worker, column, asc)
		return nil
	})
}


//
// Sheet
//

// Sheet a refrence view for table. 
// Methods should only be called by UI thread.
type Sheet struct {
	header          []string
	sorted          []int64
	idToRow         map[int64]int
	OnHeaderChanged func()
	OnSorted        func()
}

func newSheet(eb *event.EventBus, header []string) *Sheet {
	s := &Sheet{
		header: header,
		sorted: make([]int64, 0),
		OnSorted: func() {},
		OnHeaderChanged: func() {},
	}
	eb.Subscribe(EventHiddenColumn{}, func(v event.Event) {
		e := v.(EventHiddenColumn)
		header := make([]string, 0)
		for i, label := range models.BookEntryFields() {
			if !e.Hidden[i] {
				header = append(header, label)
			}
		}
		s.header = header
		s.OnHeaderChanged()
	})
	eb.Subscribe(EventSnapshotSorted{}, func(v event.Event){
		ss := snapshot.Current.Load()
		e := v.(EventSnapshotSorted)
		if ss.Version() == e.Version {
			s.sorted = e.Sorted
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
	cb        *command.CommandBus
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
		Asc: s.Ascending,
		Column: s.Column,
	})
}

//
// Searchable
//


type Searchable struct {
	headers  []string
	OnChangedOptions  func()
	onChangedSearchBy func(s string)
}

func newSearchable(headers []string) *Searchable {
	s := &Searchable{
		headers: headers,
		OnChangedOptions: func() {},
		onChangedSearchBy: func(_ string) {},

	}
	return s
}

func (s *Searchable) SetSearchBy(h string) {
	s.onChangedSearchBy(h)
}

func (s *Searchable) setSelectable(headers []string) {
	s.headers = headers
	s.OnChangedOptions()
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
	cb       *command.CommandBus
	debounce func(func())
	column   string
}

func newSearching(column string, cb *command.CommandBus) *Searching {
	sr := &Searching{
		column: column,
		cb: cb,
		debounce: worker.Debounce(time.Millisecond * 300),
	}
	return sr
}

func (s *Searching) setSearchColumn(h string) {
	s.column = h
}

func (s *Searching) Search(pattern string) {
	s.debounce(func() {
		s.cb.Dispatch(CommandSearch{Pattern: pattern, Column: s.column})
	})
}

//
// Selected
//

// Selected 
type Selected struct {
	cb         *command.CommandBus
	OnSelected func(Point, bool)
	onSelected func(Point, bool)
}

func newSelected(cb *command.CommandBus) *Selected {
	es := &Selected{
		cb: cb,
		OnSelected: func(_ Point, _ bool) {},
	}
	es.onSelected = func(p Point, has bool) {
		es.OnSelected(p, has)
	}
	return es
}

func (es *Selected) Set(p Point, ok bool) {
	es.cb.Dispatch(CommandSheetSelect{Point: p, Has: ok})
}

//
// Search Selection
//

type SearchSelection struct {
	eb *event.EventBus
	cb *command.CommandBus
	ssVersion int64
	selection []snapshot.Point
	position  int
}

func newSearchSelection(eb *event.EventBus, cb *command.CommandBus) *SearchSelection {
	sc := &SearchSelection{
		eb: eb,
		cb: cb,
		position: -1,
	}
	eb.Subscribe(EventSnapshotSearched{}, func(v event.Event){
		e := v.(EventSnapshotSearched)
		sc.ssVersion = e.Version
		sc.selection = e.Points
		sc.position = 0
	})
	return sc
}

func (es *SearchSelection) Next() {
	if len(es.selection) == 0 {
		return
	}
	es.position += 1
	if es.position >= len(es.selection) {
		es.position = 0
	}
	es.selected()
}

func (es *SearchSelection) Prev() {
	if len(es.selection) == 0 {
		return
	}
	es.position -= 1
	if es.position < 0 {
		es.position = len(es.selection)-1
	}
	es.selected()
}

func (es *SearchSelection) selected() {
	p := es.selection[es.position]
	es.cb.Dispatch(CommandSnapshotSelect{
		Version: es.ssVersion,
		Point: p,
		Has: true,
	})
}


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
		eb: eb,
	}
	cs.eb.Subscribe(EventSnapshotSorted{}, func(v event.Event) {
		e := v.(EventSnapshotSorted)
		cfg.UI.TableSortBy = e.Column
		cfg.UI.TableAscending = e.Asc
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
	ts.eb.Notify(EventHiddenColumn{
		Hidden: hidden,
	})
}

//
// Functions and Helpers
//

func snapshotSearchWithContext(ctx context.Context, ss *snapshot.Snapshot, by string, pattern string) ([]snapshot.Point, []int, error) {
	var trv search.Traverser
	if by == ColumnAll {
		trv = newTableTraverse(ss)
	} else {
		idx := slices.Index(models.BookEntryFields(), by)
		if idx == -1 {
			return nil, nil, fmt.Errorf("search %s: invalid column label", by)
		}
		trv = newColumnTraverse(ss, idx)
	}
	srch := (&search.Searcher{}).Set(trv, pattern)

	//points, score := searchSearcherWithContext(ctx, srch)
	results := searchSearcherWithContext(ctx, srch)

	scores := make([]int, len(results))
	points := make([]snapshot.Point, len(results))
	for i, r := range results {
		if ctx.Err() != nil {
			return points, scores, nil
		}
		scores[i] = r.Score
		p := snapshot.Point{
			Row: r.Point.Row,
			Col: r.Point.Col,
		}
		points[i] = p
	}
	return points, scores, nil
}

type SearchResult struct {
	Point search.Point
	Score int
}

func searchSearcherWithContext(ctx context.Context, srch *search.Searcher) ([]SearchResult) {
	
	results := make([]SearchResult, 0)

	for srch.Next() {
		if ctx.Err() != nil {
			return []SearchResult{}
		}
		point := srch.Point()
		score := srch.Score()
		if score == -1 {
			continue
		}
		r := SearchResult{
			Score: score,
			Point: point,
		}
		results = append(results, r)
	}

	if len(results) == 0 {
		return []SearchResult{}
	}

	slices.SortFunc(results, func(a, b SearchResult) int {
		r := cmp.Compare(a.Score, b.Score)
		if r == 0 {
			return cmp.Compare(a.Point.Row, b.Point.Row)
		}
		return r * -1
	})
	return results
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

func snapshotSort(ss *snapshot.Snapshot, column string, asc bool) ([]int64, error) {
	idx := slices.Index(models.BookEntryFields(), column)
	if idx == -1 {
		return nil, fmt.Errorf("sort %s: invalid column", column)
	}
	comp, err := app.CompareBookEntry(idx, asc)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, ss.Length())
	for i := range ss.Length() {
		id, _ := ss.RowToID(i)
		ids[i] = id
	}
	slices.SortFunc(ids, func(a, b int64) int {
		ba, _ := ss.GetBookEntryByID(a)
		bb, _ := ss.GetBookEntryByID(b)
		return comp(*ba, *bb)
	})
	return  ids, nil
}

//func sortIDs(ids []int64, ss *snapshot.Snapshot, column string, asc bool) {
//	colIdx := slices.Index(models.BookEntryFields(), column)
//	comp, err := app.CompareBookEntry(colIdx, asc)
//	if err != nil {
//		log.Println("sorting:", err)
//		return
//	}
//	slices.SortFunc(ids, func(a, b int64) int {
//		bookA, _ := ss.GetBookEntryByID(a)
//		bookB, _ := ss.GetBookEntryByID(b)
//		return comp(*bookA, *bookB)
//	})
//}

//func isValidVersion(curr, result *snapshot.Snapshot) bool {
//	return curr != nil && curr.Version() != result.Version()
//}

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
		return 0, fmt.Errorf("to_sheet_column %d: index out of range", column)
	}
	label := models.BookEntryFields()[column]
	col := slices.Index(header, label)
	if col == -1 {
		return 0, fmt.Errorf("to_sheet_column %s: label not found", label)
	}
	return col, nil
}


