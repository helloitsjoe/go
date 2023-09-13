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

func TestIndex(t *testing.T) {
	e := router.New("../")
	Register(e)
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, Index(c))
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

func TestRegisterUserNoInput(t *testing.T) {
	e := router.New("../")
	Register(e)
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, RegisterUser(c))
	r := rec.Body.String()
	fmt.Println(r)
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Name and password must be provided")
}

// func TestRegisterUserValid(t *testing.T) {
// 	e := router.New("../")
// 	Register(e)
// 	req := httptest.NewRequest(echo.GET, "/", strings.NewReader("username=Alice&password=123"))
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	assert.NoError(t, RegisterUser(c))
// 	r := rec.Body.String()
// 	fmt.Println(r)
// 	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
// 	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
// 	assert.Contains(t, r, "div class=\"error\"")
// 	assert.Contains(t, r, "Name and password must be provided")
// }
