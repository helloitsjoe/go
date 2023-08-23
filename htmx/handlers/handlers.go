package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	file, err := os.ReadFile("static/index.html")
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "There was a problem")
	}
	return c.HTML(http.StatusOK, string(file))
}

func Update(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	return c.HTML(http.StatusOK, "{\"message\": \"Updated!\"}")
}
