package status

import (
	"fmt"
	"errors"
	"log/slog"
)


type StatusCode int 
const (
	LoadConfig     StatusCode = iota
	OpenedDatabase
)

type Status struct {
	Code StatusCode
	ErrLog []error
}



const (
	LevelDebug slog.Level = slog.LevelDebug
	LevelInfo  slog.Level = slog.LevelInfo
	LevelWarn  slog.Level = slog.LevelWarn
	LevelError slog.Level = slog.LevelError
)

type Kind string
const (
	FileNotFound   Kind = "file not found"
	FailedToCreate Kind = "failed to create"
	FailedToOpen   Kind = "failed to open"
	FailedToDecode Kind = "failed to decode"
	Unexpected     Kind = "unexpected error"
)

type Op string

type Error struct {
	Op   Op
	Kind Kind    
	Err  error
	Severity slog.Level
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Op, e.Kind, e.Err)
}

func E(args ...any) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case error:
			e.Err = arg
		case string:
			e.Err = errors.New(arg)
		case slog.Level:
			e.Severity = arg
		default:
			panic("invalid input to E")
		}
	}
	return e
}

