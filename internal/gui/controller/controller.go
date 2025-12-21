package controller

import (

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/config"
	storeDriver "github.com/dubbersthehoser/mayble/internal/storage/driver"
)

type Controller struct {
	App        *app.App

	List      *BookLoanList
	Searcher  *BookLoanSearcher
	Editer    *BookEditer
	Config    *config.Config

	Broker    *emiter.Broker
}

func New(cfg *config.Config) (*Controller, error) {
	var c Controller
	storage, err := storeDriver.Load(cfg.DBDriver, cfg.DBFile)
	if err != nil {
		return nil, err
	}
	a := app.New(storage)
	if err != nil {
		return nil, err
	}

	c.Broker = &emiter.Broker{}
	c.Config = cfg

	c.App = a
	c.List = NewBookLoanList(c.Broker, a)
	c.Searcher = NewBookLoanSearcher(c.Broker, &c.List.list)
	c.Editer = NewBookEditer(c.Broker, a)

	return &c, nil
}

func (c *Controller) Reset() error {
	storage, err := storeDriver.Load(c.Config.DBDriver, c.Config.DBFile)
	if err != nil {
		return err
	}
	a := app.New(storage)
	if err != nil {
		return err
	}
	c.App.Close()
	c.App = a
	c.List.SetApp(a)
	c.Editer.SetApp(a)
	return nil
}

