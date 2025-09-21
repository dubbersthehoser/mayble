package controller

type Command interface {
	Execute() error
}

type CommandBookLoanDelete struct {}

type CommandBookLoanCreate struct {}

type CommandBookLoanUpdate struct {}

func (c *Controller) GetCommand(driver string) Command {
	
}
