package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Auth middleware", c.Cookies())

		username, err := c.Cookie("username")

		if err != nil && !strings.Contains(err.Error(), "named cookie not present") {
			fmt.Println(err)
			return err
		}

		if username != nil {
			fmt.Println("username", username.Value)
			c.Set("username", username.Value)
		}

		return next(c)
	}
}
