package handlers

import (
	"encoding/json"
	"fmt"
	"htmx/router"
	"net/http/httptest"
	"net/url"
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

func TestGetUsersHtmx(t *testing.T) {
	e := router.New("../")
	Register(e)
	req := httptest.NewRequest(echo.GET, "/users", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, AllUsers(c))
	r := rec.Body.String()
	assert.Contains(t, r, "Alice")
	assert.Contains(t, r, "Bob")
	assert.Contains(t, r, "Carl")
}

func TestGetUsersJson(t *testing.T) {
	e := router.New("../")
	Register(e)
	req := httptest.NewRequest(echo.GET, "/users?format=json", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, AllUsers(c))
	r := rec.Body.Bytes()
	users := map[string]User{}
	err := json.Unmarshal(r, &users)
	assert.Nil(t, err)

	alice, bob, carl := users["Alice"], users["Bob"], users["Carl"]

	assert.Equal(t, alice.Username, "Alice")
	assert.Equal(t, bob.Username, "Bob")
	assert.Equal(t, carl.Username, "Carl")

	// Bit annoying way to check that the Ids are unique
	assert.NotEqual(t, alice.Id, bob.Id)
	assert.NotEqual(t, bob.Id, carl.Id)
	assert.NotEqual(t, alice.Id, carl.Id)
}

func TestRegisterUserNoInput(t *testing.T) {
	e := router.New("../")
	Register(e)
	req := httptest.NewRequest(echo.GET, "/register", strings.NewReader(""))
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

func TestRegisterUserValid(t *testing.T) {
	e := router.New("../")
	Register(e)
	form := url.Values{}
	form.Add("username", "New User")
	form.Add("password", "123")
	req := httptest.NewRequest(echo.POST, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, RegisterUser(c))
	r := rec.Body.String()
	fmt.Println(r)
	assert.Contains(t, r, "<span id=\"num-users\" hx-swap-oob=\"true\">4</span>")
	assert.Contains(t, r, "<span id=\"logged-in-as\" hx-swap-oob=\"true\">Logged in as New User</span>")
	assert.Contains(t, r, "New User, you're in.")
}

// func TestLogin(t *testing.T) {
// 	e := router.New("../")
// 	Register(e)
// 	form := url.Values{}
// 	form.Add("username", "Alice")
// 	form.Add("password", "123")
// 	req := httptest.NewRequest(echo.POST, "/register", strings.NewReader(form.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	assert.NoError(t, RegisterUser(c))
// 	r := rec.Body.String()
// 	fmt.Println(r)
// 	assert.Contains(t, r, "<span id=\"num-users\" hx-swap-oob=\"true\">4</span>")
// 	assert.Contains(t, r, "<span id=\"logged-in-as\" hx-swap-oob=\"true\">Logged in as New User</span>")
// 	assert.Contains(t, r, "New User, you're in.")
// }
