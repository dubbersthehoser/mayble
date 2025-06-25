package app

import (
	"io"
	"log"
	"os"
	_ "fmt"
)

type Logs struct {
	writers []io.Writer
}
func (l *Logs) Write(p []byte) (int, error) {
	for _, w := range l.writers {
		_, err := w.Write(p)
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}
func (l *Logs) Register(w io.Writer) {
	l.writers = append(l.writers, w)
}
func NewLogs() *Logs{
	l := &Logs{
		writers: []io.Writer{os.Stderr},
	}
	log.SetOutput(l)
	return l
}
