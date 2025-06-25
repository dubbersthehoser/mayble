package commands

import (
	"github.com/dubbersthehoser/mayble/internal/app"
)

type Command {
	Name  string
	args []string
}

type Commands struct {
	Table map[string]func(*app.State, *Command) error
	Help  map[string]string
}
func newCommands() *Commands{
	c := &Commands{
		Table: make(map[string]func(*app.State, Command) error),
		Help: make(map[string]string),
	}
	return c
}
func (c *Commands) Register(name string, handler func(*app.State, *Command) error, help string) {
	c.Table[name] = handler
	c.Help[name] = help
}
func (c *Commands) Run(s *app.State, cmd *Commands) error {
	handler, ok := c.Table[cmd.Name]
	if !ok {
		return fmt.Errorf("'%s' command not found", cmd.Name)
	}
	return handler(s, cmd)
}
