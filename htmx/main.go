package main

import (
	"htmx/db"
	"htmx/handlers"
	"htmx/router"
	"htmx/user"
)

func main() {
	e := router.New("")
	d := db.CreateDB()
	user.SeedUsers(d)
	handlers.Register(e, d)
	e.Logger.Fatal(e.Start(":8080"))
}
