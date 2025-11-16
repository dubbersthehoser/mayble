package command


type Command interface {
	Do(storage.Storage)   error
	Undo(storage.Storage) error
}

type Stack struct {
	items []Command
}
func NewStack() *CommandStack {
	c := CommandStack{
		items: make([]Command, 0),
	}
	return &c
}

func (s *Stack) Pop() Command {
	length := len(s.items)
	if length == 0 {
		return nil
	}
	cmd := s.items[length-1]
	cs.items = s.items[:length-1]
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


