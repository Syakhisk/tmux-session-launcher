package action

import (
	"io"
	"os"
)

type Action struct {
}

func New() *Action {
	return &Action{}
}

func (a *Action) Next() error {
	f := os.NewFile(uintptr(3), "pipe")
	defer f.Close()
	io.WriteString(f, "hello from child\n")
	return nil
}
