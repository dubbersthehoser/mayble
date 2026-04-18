package viewmodel

import (
	"errors"
	"os"
	"testing"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/app"
)

func TestMenuUI(t *testing.T) {
	b := &bus.Bus{}
	db, err := database.OpenMem()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Conn.Close()
	cfg := &config.Config{}
	as := app.NewService(cfg, db)
	err = db.Conn.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	DBPath := binding.NewString()

	menu := NewMenuVM(b, as, as, DBPath)
	_ = menu

	t.Run("DatabaseCreateAndOpen", func(t *testing.T) {
		testDatabaseCreateAndOpen(t, menu)
	})
}

func testDatabaseCreateAndOpen(t *testing.T, menu *MenuVM) {
	file, err := os.CreateTemp("", "*.db")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	dbPath := file.Name()
	file.Close()

	menu.bus.Subscribe(bus.Handler{
		Name: msgUserInfo,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_info:", e.Data.(string))
		},
	})
	menu.bus.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_error:", e.Data.(string))
		},
	})
	menu.bus.Subscribe(bus.Handler{
		Name: msgUserSuccess,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_success:", e.Data.(string))
		},
	})

	t.Run("create", func(t *testing.T) {
		test_createDatabase(t, menu, dbPath)
	})
}

func test_createDatabase(t *testing.T, menu *MenuVM, path string) {
	var err error

	menu.CreateDatabase(path, err)

	file, err := os.CreateTemp("", "")
	path = file.Name()
	file.Close()
	menu.CreateDatabase(path, err)

	expect := path + ".db"

	_, err = os.Lstat(expect)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	menu.CreateDatabase("", nil)

	var ok bool
	err = errors.New("invalid permissions")
	id := menu.bus.Subscribe(busMsgTestHelper(t, msgUserError, func(s string) {
		ok = true
		expect := "invalid permissions"
		if expect != s {
			t.Fatalf("expect message '%s', got '%s'", expect, s)
		}
	}))
	menu.CreateDatabase("", err)
	if !ok {
		t.Fatalf("expected message")
	}
	menu.bus.Unsubscribe(id)
}

