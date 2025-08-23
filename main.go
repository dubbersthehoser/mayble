package main

import (
	"embed"

	//"github.com/dubbersthehoser/mayble/internal/app"
	//"github.com/dubbersthehoser/mayble/internal/database"
	  gui "github.com/dubbersthehoser/mayble/internal/gui/"
)

//go:embed sql/schemas
var schemaFS  embed.FS
var schemaDir string = "sql/schemas"

func main() {
	gui.Run()
}
