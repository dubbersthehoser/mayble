package viewmodel

import (
	"time"
	"log"
	"errors"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/app"
)

const (
	msgUserError   string = "message.user.error"
	msgUserSuccess string = "message.user.success"
	msgUserInfo    string = "message.user.info"
)

var (
	errDBFile error = errors.New("database file error")
)

const (
	BodyData int = iota
	BodyForm
	BodyMenu
)

type subjectRetriever interface {
	AddListener(func())
	repo.BookRetriever
}

type subjectGenreRetriever interface {
	AddListener(func())
	repo.GenreRetriever
}

type databaseOpener interface {
	OpenDatabase(s string) error
}

type MainUI struct {
	bus     *bus.Bus
	cfg     *config.Config
	service *app.Service

	genres      *UniqueGenres
	store       repo.BookStore

	OpenedBody  binding.Int
	DBFile      binding.String
	HasDatabase binding.Bool

	Error   binding.String
	Success binding.String
	Info    binding.String
	Clear   binding.Bool
}

func NewMainUI(cfg *config.Config) *MainUI {

	b := &bus.Bus{}
	as := app.NewService(cfg)

	hasDatabase := true
	if err := as.OpenDatabase(cfg.DBFile); err != nil {
		log.Println("ERROR:", err.Error())
		hasDatabase = false
	}


	var store repo.BookStore = newStoreUserMessaging(as, b)
	mu := &MainUI{
		bus:        b,
		cfg:        cfg,
		service:    as,

		store:     store,
		genres:    NewUniqueGenres(b, as),

		OpenedBody:  binding.NewInt(),
		DBFile:      binding.NewString(),
		HasDatabase: binding.NewBool(),

		Error:   binding.NewString(),
		Success: binding.NewString(),
		Info:    binding.NewString(),
		Clear:   binding.NewBool(),
	}

	_ = mu.HasDatabase.Set(hasDatabase)

	mu.SetOpenBody(mu.GetOpenBody())
	_ = mu.DBFile.Set(cfg.DBFile)

	// TODO Refactor this out into its own type.
	//
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

func (m *MainUI) SetOpenBody(w int) {
	_ = m.OpenedBody.Set(w)
	m.cfg.UI.OpenBody = w
}

func (m *MainUI) GetOpenBody() int {
	return m.cfg.UI.OpenBody
}

func (m *MainUI) GetMenu() *Menu {
	return NewMenu(m.bus, m.service, m.service, m.DBFile)
}

func (m *MainUI) GetTable() *Table {
	return NewTable(m.bus, &TableConfig{cfg: m.cfg}, m.store, m.service, m.genres)
}

func (m *MainUI) GetBookSubmissionForm() *BookSubmissionForm {
	return NewBookSubmissionForm(m.bus, m.store, m.genres)
}

func (m *MainUI) GetBodySwitcher() *BodySwitcher {
	return newBodySwitcher(m.cfg, m.HasDatabase)
}


type StatusLine struct {
}


type BodyButton struct {
	ID        int
	OnLock    func()
	OnUnlock  func()
	Window BodyWindow
}

type BodyWindow struct {
	OnShow func()
	OnHide func()
}

type BodySwitcher struct {
	cfg     *config.Config
	Buttons map[int]BodyButton
	hasDatabase binding.Bool
}

func newBodySwitcher(cfg *config.Config, hasDatabase binding.Bool) *BodySwitcher {
	bs := &BodySwitcher{
		cfg: cfg,
		hasDatabase: hasDatabase,
	}

	hasDatabase.AddListener(binding.NewDataListener(func() {
		ok, _ := hasDatabase.Get()
		bs.setDatabaseState(ok)
	}))

	return bs
}

func (bs *BodySwitcher) setDatabaseState(ok bool) {
	if !ok {
		bs.Switch(BodyMenu)
		bs.Buttons[BodyForm].OnLock()
		bs.Buttons[BodyData].OnLock()
	} else {
		bs.Buttons[BodyForm].OnUnlock()
		bs.Buttons[BodyData].OnUnlock()
	}
}

func (bs *BodySwitcher) SetBodies(buttons ...BodyButton) {
	bs.Buttons = make(map[int]BodyButton)
	for _, btn := range buttons {
		bs.Buttons[btn.ID] = btn
	}
	ok, _ := bs.hasDatabase.Get()
	bs.setDatabaseState(ok)
}

func (bs *BodySwitcher) Switch(id int) {
	for _, b := range bs.Buttons {
		if b.ID == id {
			b.OnLock()
			b.Window.OnShow()
			bs.cfg.UI.OpenBody = b.ID
		} else {
			b.OnUnlock()
			b.Window.OnHide()
		}
	}
}
