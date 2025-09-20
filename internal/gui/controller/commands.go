package controller

type Command interface {
	Execute() error
}

type CommandBookLoanDelete struct {}

type CommandBookLoanCreate struct {}

type CommandBookLoanUpdate struct {}

type CommandUndo struct {}

type CommandRedo struct {}
