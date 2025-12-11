package main

import (
	"os"
	"log"

	"github.com/dubbersthehoser/mayble/internal/launcher"
)


func main() {
	if err := launcher.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
