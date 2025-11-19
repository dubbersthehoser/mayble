package stub

type Command struct {
	Count int
	Label string
}

func (c *Command) Do() error {
	c.Count += 1
	return nil
}

func (c *Command) Undo() error {
	c.Count -= 1
	return nil
}
