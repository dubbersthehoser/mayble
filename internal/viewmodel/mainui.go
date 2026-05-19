package viewmodel

import (
	"time"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/app"
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
	cfg     UIConfig
	errList []error

	genres      *UniqueGenres
	store       repo.BookStore
	retriever   repo.BookRetriever
	dbOpener    DatabaseOpener
	fileHandler repo.CSVHandler

	OpenedBody binding.Int
	DBFile     binding.String

	Error   binding.String
	Success binding.String
	Info    binding.String
	Clear   binding.Bool
}

func NewMainUI(cfg *config.Config, db *database.Database, errs []error) *MainUI {

	b := &bus.Bus{}
	as := app.NewService(cfg, db)
	var store repo.BookStore = newStoreUserMessaging(as, b)
	store = newStoreNotifyChanged(store, b)
	mu := &MainUI{
		OpenedBody: binding.NewInt(),
		bus:        b,
		cfg:        &cfg.UI,

		store:     store,
		genres:    NewUniqueGenres(b, as),
		retriever: as,

		DBFile:   binding.NewString(),
		dbOpener: as,
		fileHandler: as,

		errList: errs,


		Error:   binding.NewString(),
		Success: binding.NewString(),
		Info:    binding.NewString(),
		Clear:   binding.NewBool(),
	}

	mu.SetBody(mu.GetBody())

	mu.DBFile.Set(cfg.DBFile)

	// to clear info line
	countDown := time.Duration(time.Minute / 10)
	timer := time.NewTimer(0)
	clearLine := func() {
		go func() {
			_ = mu.Clear.Set(false)
			timer.Stop()
			timer.Reset(countDown)
			<-timer.C
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
			_ = mu.Success.Set("")
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
			_ = mu.Success.Set("")
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

func (m *MainUI) SetBody(w int) {
	_ = m.OpenedBody.Set(w)
	m.cfg.SetWindowBody(w)
}

func (m *MainUI) GetBody() int {
	return m.cfg.GetWindowBody()
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
	return NewMenuVM(m.bus, m.fileHandler, m.dbOpener, m.DBFile)
}

func (m *MainUI) GetTableVM() *TableVM {
	return NewTableVM(m.bus, m.cfg, m.retriever)
}

func (m *MainUI) GetTableControllersVM() *TableControllersVM {
	return NewTableControllersVM(m.bus, m.retriever, m.store, m.genres)
}

func (m *MainUI) GetBookSubmissionForm() *BookSubmissionForm {
	return NewBookSubmissionForm(m.bus, m.store, m.genres)
}

