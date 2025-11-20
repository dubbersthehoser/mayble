package command

type Command interface {
	Do() error
	Undo() error
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


type Queue struct {
	items []Command
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]Command, 0),
	}
}

func (q *Queue) Enqueue(cmd Command) {
	q.items = append(q.items, cmd)
}

func (q *Queue) Dequeue() Command {
	if len(q.items) == 0 {
		return nil
	}
	cmd := q.items[0]
	q.items = q.items[1:]
	return cmd
}

func (q *Queue) Length() int {
	return len(q.items)
}

func (q *Queue) Clear() {
	q.items = make([]Command, 0)
}
