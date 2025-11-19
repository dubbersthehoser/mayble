package app


/**********************
	Commands
***********************/

/* Import */

// NOTE need to re-think import system implementation.
type commandImportBookLoans struct {
	addedIDs  []int64
	bookLoans []data.BookLoan
}
func (c *commandImportBookLoans) Do(s storage.BookLoanStore) error {
	c.addedIDs = make([]int64, len(c.bookLoans))
	for i, BookLoan := range c.bookLoans {
		id, err := createBookLoan(s, BookLoan)
		if err != nil {
			return fmt.Errorf("app: import: %w", err)
		}
		c.addedIDs[i] = BookLoan
	}
	return nil
}
func (c *commandImportBookLoans) Undo(s storage.BookLoanStore) error {
	for i, id := range c.addedIDs {
		book := c.bookLoans[i]
		book.ID = id
		err = deleteBookLoan(s, &book)
		if err != nil {
			return err
		}
	}
	return nil
}


/* Create */

type commandCreateBookLoan struct {
	bookLoan *BookLoan
}

func (c *commandCreateBookLoan) Do(s storage.Storage) error {
	_, err := createBookLoan(s, c.bookLoan)
	return err
}

func (c *commandCreateBookLoan) Undo(s storage.Storage) error {
	return deleteBookLoan(s, c.bookLoan)
}


/* Delete */

type commandDeleteBookLoan struct {
	bookLoan *data.BookLoan
}

func (c *commandDeleteBookLoan) Do(s storage.Storage) error {
	return deleteBookLoan(s, c.bookLoan)
}

func (c *commandDeleteBookLoan) Undo(s storage.Storage) error {
	_, err := createBookLoan(s, c.bookLoan)
	return err
}


/* Update */

type commandUpdateBookLoan struct {
	bookLoan *BookLoan
	prevBookLoan *BookLoan
}

func (c *commandUpdateBookLoan) Do(s storage.Storage) error {
	bookLoan, err := getBookLoanByID(s, c.bookLoan.ID)
	if err != nil {
		return err
	}
	if c.prevBookLoan == nil {
		c.prevBookLoan = &bookLoan
	} else {
		bookLoan = *c.prevBookLoan
		c.prevBookLoan = c.bookLoan
		c.bookLoan = &bookLoan
	}
	return updateBookLoan(s, c.bookLoan)
}

func (c *commandUpdateBookLoan) Undo(s storage.Storage) error {
	book := c.prevBookLoan
	c.prevBookLoan = c.bookLoan
	c.bookLoan = book
	return s.UpdateBookLoan(c.bookLoan)
}





/****************************************
        Command Storage Manager
*****************************************/

type manager struct {
	store storage.Storage
	undos *command.Stack
	redos *command.Stack
	queue []command.Command
}

func newManager(store storage.Storage) *manager{
	m := manager{
		undos: command.NewStack(),
		redos: command.NewStack(),
		store: store,
	}
	return &m
}

func (m *manager) execute(cmd command.Command) error {
	if err := cmd.Do(m.store); err != nil {
		return err
	}
	m.undos.Push(cmd)
	m.redos.Clear()
	return nil
}

func (m *manager) unExecute() error {
	cmd := m.undos.Pop()
	if cmd == nil {
		return nil
	}
	err := cmd.Undo(m.store)
	if err != nil {
		return err
	}
	m.redos.Push(cmd)
	return nil
} 

func (m *manager) reExecute() error {
	cmd := m.redos.Pop()
	if cmd == nil {
		return nil
	}
	err := cmd.Do(m.store)
	if err != nil {
		return err
	}
	m.undos.Push(cmd)
	return nil
}

func (m *manager) enqueue(cmd command.Command) {
	m.queue = append(m.queue, cmd)
}

func (m *manager) dequeue() command.Command {
	if len(m.queue) == 0 {
		return nil
	}
	cmd := m.queue[0]
	m.queue = m.queue[1:]
	return cmd
}
