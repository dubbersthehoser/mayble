package cli

import (
	"os"

	"github.com/dubbersthehoser/mayble/internal/broker"
)

type Command struct {
	Handler func() error
}
func (c *Command) Run() error {
	return c.Handler()
}

type CommandLookup struct {
	Handlers map[string]*Command
}
func NewCommandLookup() *CommandLookup {
	return &CommandLookup{
		Handlers: make(map[string]Command),
	}
}
func(cl *CommandLookup) Register(name string, c *Command) {
	cl.Handlers[name] = c
}
func (cl *CommandLookup) Lookup(name string) (*Command, error) {
	cmd, ok := cl.Handlers[namel]
	if !ok {
		return nil, errors.New("command not found")
	}
	return cmd, nil
}

func AddBookHander(args []string, b *broker.Broker) *Command {
	return &Command{
		Args: args,
		Handler: func() error {
			fmt.Println("add command ran")
			_, := b.Emit("book.remove", nil)
			return nil
		}
	}
}
func AddBookHander(args []string, b *broker.Broker) *Command {
	return &Command{
		Args: args,
		Handler: func() error{
			fmt.Println("add command ran")
			_, := b.Emit("book.add", nil)
			return nil
		}
	}
}

func Run() error {
	
	b := broker.NewBroker()

	cl := NewCommandLookup()
	cl.Register("add", AddBookHandler(os.Args, b))
	cl.Register("remove", RemoveBookHandler(os.Args, b))
	
}
