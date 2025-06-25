package main

import (
	"embed"
	"time"
	"context"
	"log"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/database"
)

//go:embed sql/schemas
var schemaFS  embed.FS
var schemaDir string = "sql/schemas"

func main() {

	var err error
	app.SchemaFS  = schemaFS
	app.SchemaDir = schemaDir
	state := app.Init()
	ctx := context.Background()

	prams := database.CreateBookParams{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Title: "Example Title",
		Author: "Example Author",
		Genre: "Sampling",
		Ratting: 5,
	}

	_, err = state.DB.Queries.CreateBook(ctx, prams)
	if err != nil {
		log.Fatal(err)
	}

	books, err := state.DB.Queries.GetAllBooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range books {
		println(book.Title)
	}

}
