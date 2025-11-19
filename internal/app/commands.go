package app

/**********************
	Commands
***********************/

/* Import */

// NOTE need to re-think import system implementation.

type commandImportBookLoans struct {
	store    storage.BookLoanStore
	addedIDs  []int64
	bookLoans []BookLoan
}
func newCommandImportBookLoans(books []BookLoan) func(storage.BookLoanStore) *commandImportBookLoans {
	return func(s storage.BookLoanStore) {
		return &commandCreateBookLoan{
			bookLoans: book,
			store: s,
		}
	}
}

func (c *commandImportBookLoans) Do() error {
	c.addedIDs = make([]int64, len(c.bookLoans))
	for i, BookLoan := range c.bookLoans {
		id, err := createBookLoan(c.store, BookLoan)
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
		err = deleteBookLoan(c.store, &book)
		if err != nil {
			return err
		}
	}
	return nil
}


/* Create */

type commandCreateBookLoan struct {
	store    storage.BookLoanStore
	bookLoan *BookLoan
}
func newCommandCreateBookLoan(book *BookLoan) func(storage.BookLoanStore) *commandCreateBookLoan {
	return func(s storage.BookLoanStore) {
		return &commandCreateBookLoan{
			bookLoan: book,
			store: s,
		}
	}
}

func (c *commandCreateBookLoan) Do() error {
	_, err := createBookLoan(c.store, c.bookLoan)
	return err
}

func (c *commandCreateBookLoan) Undo() error {
	return deleteBookLoan(c.store, c.bookLoan)
}


/* Delete */

type commandDeleteBookLoan struct {
	store    storage.BookLoanStore
	bookLoan *data.BookLoan
}
func newCommandDeleteBookLoan(book *BookLoan) func(storage.BookLoanStore) *commandDeleteBookLoan {
	return func(s storage.BookLoanStore) {
		return &commandDeleteBookLoan{
			bookLoan: book,
			store: s,
		}
	}
}

func (c *commandDeleteBookLoan) Do() error {
	return deleteBookLoan(c.store, c.bookLoan)
}

func (c *commandDeleteBookLoan) Undo() error {
	_, err := createBookLoan(c.store, c.bookLoan)
	return err
}


/* Update */

type commandUpdateBookLoan struct {
	store    storage.BookLoanStore
	bookLoan *BookLoan
	prevBookLoan *BookLoan
}

func newCommandUpdateBookLoan(book *BookLoan) func(storage.BookLoanStore) commandUpdateBookLoan {
	return func(s storage.BookLoanStore) *commandUpdateBookLoan {
		return &commandUpdateBookLoan{
			bookLoan: book,
			store: s,
		}
	}
}

func (c *commandUpdateBookLoan) Do() error {
	bookLoan, err := getBookLoanByID(c.store, c.bookLoan.ID)
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
	return updateBookLoan(c.store, c.bookLoan)
}

func (c *commandUpdateBookLoan) Undo() error {
	book := c.prevBookLoan
	c.prevBookLoan = c.bookLoan
	c.bookLoan = book
	return updateBookLoan(c.store, c.bookLoan)
}





/****************************************
        Command Storage Manager
*****************************************/

type manager struct {
	store storage.BookLoanStore
	undos *command.Stack
	redos *command.Stack
	queue []command.Command
}

func newManager(store storage.BookLoanStore) *manager{
	m := manager{
		undos: command.NewStack(),
		redos: command.NewStack(),
		store: store,
	}
	return &m
}

// execute command
func (m *manager) execute(cmdSetup func(storage.BookLoanStore) command.Command) error {
	cmd := cmdSetup(m.store)
	if err := cmd.Do(); err != nil {
		return err
	}
	m.undos.Push(cmd)
	m.redos.Clear()
	return nil
}

// unExecute command
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

// reExecute an undo'ed command
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

// enqueue command into queue
func (m *manager) enqueue(cmdSetup func(storage.BookLoanStore) command.Command) {
	cmd := cmdSetup(m.store)
	m.queue = append(m.queue, cmd)
}

// dequeue command out of queue
func (m *manager) dequeue() command.Command {
	if len(m.queue) == 0 {
		return nil
	}
	cmd := m.queue[0]
	m.queue = m.queue[1:]
	return cmd
}
