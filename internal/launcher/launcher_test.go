package launcher

import (
	"os"
	"testing"
	"syscall"
)

func TestGetDBFile(t *testing.T) {
	
	println("hello!")
	root, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	println("hello!")
	err = syscall.Chroot(root)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
