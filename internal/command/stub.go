package command

type StubCommand struct {
	Count int
	Label string
}

func (c *StubCommand) Do() error {
	c.Count += 1
	return nil
}

func (c *StubCommand) Undo() error {
	c.Count -= 1
	return nil
}

