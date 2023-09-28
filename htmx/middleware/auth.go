package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Auth middleware", c.Cookies())

		user, err := c.Cookie("userid")

		if err != nil && !strings.Contains(err.Error(), "named cookie not present") {
			fmt.Println(err)
			return err
		}

		if user != nil {
			fmt.Println("user", user.Value)
			c.Set("user", "Alice")
		}

		return next(c)
	}
}
