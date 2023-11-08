package middleware

import (
	"fmt"
	"htmx/db"
	"htmx/types"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func getLoggedInUser(id string, db db.DB) *types.User {
	// TODO: What if the user is not found?
	if id == "" {
		return nil
	}

	// Length in bytes
	if len(id) != 36 {
		return nil
	}

	validUser, _ := db.FindUser(id)

	return validUser
}

func Auth(next echo.HandlerFunc, db db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Cookies in auth middleware:", c.Cookies())

		uuid, err := c.Cookie("uuid")

		if err != nil && !strings.Contains(err.Error(), "named cookie not present") {
			fmt.Println("Error reading cookie:", err)
			return err
		}

		if uuid != nil {
			user := getLoggedInUser(uuid.Value, db)
			if user != nil {
				fmt.Println("Logged in user")
				c.Set("user", user)
			} else {
				fmt.Println("Not a logged in user")
				c.SetCookie(&http.Cookie{Name: "uuid", Value: "", HttpOnly: true, MaxAge: -1})
			}
		}

		return next(c)
	}
}
