package middleware

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Auth middleware", c.Cookies())

		uuid, err := c.Cookie("uuid")

		if err != nil && !strings.Contains(err.Error(), "named cookie not present") {
			fmt.Println("Error reading cookie:", err)
			return err
		}

		if uuid != nil {
			fmt.Println("uuid", uuid.Value)
			c.Set("uuid", uuid.Value)
		}

		return next(c)
	}
}
