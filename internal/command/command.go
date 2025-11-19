package command

import (
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/command/stub"
)

// TODO remove the storage depenency
//
//type Command interface {
//	Do() error
//	Undo() error
//}

type Command interface {
	Do(storage.BookLoanStore) error
	Undo(storage.BookLoanStore) error
}

type Stack struct {
	items []Command
}

func NewStack() *Stack {
	c := &Stack{
		items: make([]Command, 0),
	}
	return c
}

func (s *Stack) Pop() Command {
	length := len(s.items)
	if length == 0 {
		return nil
	}
	cmd := s.items[length-1]
	s.items = s.items[:length-1]
	return cmd
}

func (s *Stack) Push(cmd Command) {
	s.items = append(s.items, cmd)
}

func (s *Stack) Length() int {
	return len(s.items)
}

func (s *Stack) Clear() {
	s.items = make([]Command, 0)
}


