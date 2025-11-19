package stub

import (
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type Command struct {
	Count int
}


func (c *Command) Do(s storage.BookLoanStore) error {
}
