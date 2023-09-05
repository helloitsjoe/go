package main

import (
	"htmx/handlers"
	"htmx/router"
)

func main() {
	// d := db.CreateDB()
	// h := handlers.New(d)
	e := router.New("")
	handlers.Register(e)
	e.Logger.Fatal(e.Start(":8080"))
}
