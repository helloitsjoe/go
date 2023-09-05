package handlers

import (
	"fmt"
	"htmx/router"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	// h := New("")
	e := router.New("../")
	Register(e)
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, Index(c))
	r := rec.Body
	fmt.Println(r)
	assert.Equal(t, "foo", "bar")
	// assert.NoError(t,
}
