package viewmodel

import (
	"time"
	"slices"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/config"
)

const (
	msgUserError   string = "message.user.error"
	msgUserSuccess string = "message.user.success"
	msgUserInfo    string = "message.user.info"
	msgDataChanged string = "message.data.changed"
)


const (
	BodyData int = iota
	BodyForm
	BodyMenu
)

type MainUI struct {	

	bus     *bus.Bus
	app     *appService
	errList []error

	OpenedBody binding.Int
	DBFile     binding.String

	Error      binding.String
	Success    binding.String
	Info       binding.String
	Clear      binding.Bool
}

func NewMainUI(cfg *config.Config, db *database.Database, errs []error) *MainUI {

	b := &bus.Bus{}
	as := newAppService(b, cfg, db)
	mu := &MainUI{
		OpenedBody: binding.NewInt(),
		bus: b,
		app: as,

		errList: errs,

		DBFile: binding.NewString(),

		Error: binding.NewString(),
		Success: binding.NewString(),
		Info: binding.NewString(),
		Clear: binding.NewBool(),
	}
	_ = mu.DBFile.Set(cfg.DBFile)
	mu.DBFile.AddListener(binding.NewDataListener(func() {
		path, _ := mu.DBFile.Get()
		cfg.DBFile = path
	}))

	// to clear info line
	countDown := time.Duration(time.Minute / 10)
	timer := time.NewTimer(0)
	clearLine := func() {
		go func() {
			_ = mu.Clear.Set(false)
			timer.Stop()
			timer.Reset(countDown)
			<- timer.C
			_ = mu.Clear.Set(true)
		}()
	}

	mu.bus.Subscribe(bus.Handler{
		Name: msgUserInfo,
		Handler: func(e *bus.Event) {
			if e.Data == nil {
				return
			}
			v, ok := e.Data.(string)
			if !ok {
				return
			}
			_ = mu.Error.Set("")
			_= mu.Success.Set("")
			_ = mu.Info.Set(v)
			clearLine()
		},
	})
	mu.bus.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			if e.Data == nil {
				return
			}
			v, ok := e.Data.(string)
			if !ok {
				return
			}
			_= mu.Success.Set("")
			_ = mu.Info.Set("")
			_ = mu.Error.Set(v)
			clearLine()
		},
	})
	mu.bus.Subscribe(bus.Handler{
		Name: msgUserSuccess,
		Handler: func(e *bus.Event) {
			if e.Data == nil {
				return
			}
			v, ok := e.Data.(string)
			if !ok {
				return
			}
			_ = mu.Info.Set("")
			_ = mu.Error.Set("")
			_ = mu.Success.Set(v)
			clearLine()
		},
	})
	return mu
}


func (m *MainUI) HasErrored() bool {
	return len(m.errList) > 0
}

func (m *MainUI) Errors() []string {
	es := make([]string, len(m.errList))
	for i, e := range m.errList {
		es[i] = e.Error()
	}
	return es
}

func (m *MainUI) GetMenuVM() *MenuVM {
	return NewMenuVM(m.bus, m.app, m.DBFile)
}


func (m *MainUI) GetTableVM() *TableVM {
	return NewTableVM(m.bus, m.app)
}

func (m *MainUI) GetTableControllersVM() *TableControllersVM {
	return NewTableControllersVM(m.bus, m.app)
}

func (m *MainUI) GetCreateBookFormVM() *CreateBookForm {
	return NewCreateBookForm(m.bus, m.app)
}

type BookVM struct {
	id int64
	Title binding.String
	Author binding.String
	Genre binding.String
}
func NewBookVM(id int64, title, author, genre string) *BookVM {
	vm := &BookVM{
		id: id,
		Title: binding.NewString(),
		Author: binding.NewString(),
		Genre: binding.NewString(),
	}
	_ = vm.Title.Set(title)
	_ = vm.Author.Set(author)
	_ = vm.Genre.Set(genre)
	return vm
}


const dateFormat = "02/01/2006"

func formatDate(t *time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(dateFormat)
}

func parseDate(t string) (*time.Time, error) {
	ret, err := time.Parse(dateFormat, t)
	return &ret, err
}


const capRating = 6

func formatRating(r int) string {
	switch r {
	case 0:
		return ""
	case 1:
		return "⭐"
	case 2:
		return "⭐⭐"
	case 3:
		return "⭐⭐⭐"
	case 4:
		return "⭐⭐⭐⭐"
	case 5:
		return "⭐⭐⭐⭐⭐"
	default:
		return "ERROR"
	}
}

func Ratings() []string{
	r := make([]string, capRating)
	for i := range capRating {
		r[i] = formatRating(i)
	}
	return r
}

type listener struct {
	listeners []binding.DataListener
}

func (t *listener) notify() {
	for _, listener := range t.listeners {
		listener.DataChanged()
	}
}

func (t *listener) AddListener(l binding.DataListener) {
	if t.listeners == nil {
		t.listeners = make([]binding.DataListener, 0)
	}
	t.listeners = append(t.listeners, l)
}

func (t *listener) RemoveListener(l binding.DataListener) {
	if t.listeners == nil {
		return
	}
	t.listeners = slices.DeleteFunc(t.listeners, func(ll binding.DataListener) bool {
		return l == ll
	})
}












