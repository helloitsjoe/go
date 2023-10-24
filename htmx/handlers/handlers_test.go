package handlers

import (
	"encoding/json"
	"htmx/db"
	"htmx/router"
	"htmx/types"
	"htmx/user"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	h := NewHandlers(db.CreateDB())
	e := router.New("../")
	req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Index(c))
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

func TestGetUsersHtmx(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	req := httptest.NewRequest(echo.GET, "/users", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.AllUsers(c))
	r := rec.Body.String()
	assert.Contains(t, r, "Alice")
	assert.Contains(t, r, "Bob")
	assert.Contains(t, r, "Carl")
}

func TestGetUsersJson(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	req := httptest.NewRequest(echo.GET, "/users?format=json", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.AllUsers(c))
	r := rec.Body.Bytes()
	users := []types.User{}
	err := json.Unmarshal(r, &users)
	assert.Nil(t, err)

	alice, bob, carl := users[0], users[1], users[2]

	assert.Equal(t, alice.Username, "Alice")
	assert.Equal(t, bob.Username, "Bob")
	assert.Equal(t, carl.Username, "Carl")

	// Bit annoying way to check that the Ids are unique
	assert.NotEqual(t, alice.UUID, bob.UUID)
	assert.NotEqual(t, bob.UUID, carl.UUID)
	assert.NotEqual(t, alice.UUID, carl.UUID)
}

func TestRegisterUserNoInput(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	req := httptest.NewRequest(echo.GET, "/register", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.RegisterUser(c))
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Name and password must be provided")
}

func TestRegisterUserExistsFail(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	form := url.Values{}
	form.Add("username", "Alice")
	form.Add("password", "bar")
	req := httptest.NewRequest(echo.POST, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.RegisterUser(c))
	r := rec.Body.String()

	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "That username is already taken")

	cookie := rec.Header().Get("Set-Cookie")
	assert.NotContains(t, cookie, `uuid=`)

	assert.NotContains(t, r, "Alice, you're in.")
	assert.NotContains(t, r, "Log out")
}

func TestRegisterUserValid(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	form := url.Values{}
	form.Add("username", "New User")
	form.Add("password", "123")
	req := httptest.NewRequest(echo.POST, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.RegisterUser(c))
	r := rec.Body.String()
	cookie := rec.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, `uuid=`)
	assert.Contains(t, r, "<span id=\"num-users\" hx-swap-oob=\"true\">4</span>")
	assert.Contains(t, r, `
<span id="logged-in-as" hx-swap-oob="true">
  <span>Logged in as New User</span>
  <button hx-post="/logout" class="text-btn">Log out</button>
</span>`)
	assert.Contains(t, r, "New User, you're in.")
	assert.Contains(t, r, "Logged in as New User")
	assert.Contains(t, r, "Log out")
}

func TestLogin(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	form := url.Values{}
	form.Add("username", "Alice")
	form.Add("password", "bar")
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Login(c))
	r := rec.Body.String()
	cookie := rec.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, "uuid=")

	assert.Contains(t, r, "<span id=\"num-users\" hx-swap-oob=\"true\">3</span>")
	assert.Contains(t, r, `
<span id="logged-in-as" hx-swap-oob="true">
  <span>Logged in as Alice</span>
  <button hx-post="/logout" class="text-btn">Log out</button>
</span>`)
	// Clear error on a successful login
	assert.Contains(t, r, "<div id=\"error\" hx-swap-oob=\"true\"></div>")
	assert.Contains(t, r, "Alice, you're in.")
	assert.Contains(t, r, "Logged in as Alice")
	assert.Contains(t, r, "Log out")
}

// TODO: convert some of these to user unit tests
func TestLoginError(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	form := url.Values{}
	form.Add("username", "Alice")
	form.Add("password", "wrong-password")
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Login(c))
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Incorrect login credentials")
}

func TestLoginNoUsernamePassword(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	form := url.Values{}
	form.Add("username", "")
	form.Add("password", "")
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Login(c))
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Name and password must be provided")
}

func TestLoginNoAccount(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	form := url.Values{}
	form.Add("username", "Not a user")
	form.Add("password", "bar")
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Login(c))
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Incorrect login credentials")
}

func TestLogout(t *testing.T) {
	d := db.CreateDB()
	h := NewHandlers(d)
	user.SeedUsers(d)
	e := router.New("../")
	req := httptest.NewRequest(echo.POST, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Logout(c))
	r := rec.Body.String()
	assert.Contains(t, rec.Header().Get("Set-Cookie"), "uuid=; Max-Age=0; HttpOnly")
	assert.Contains(t, r, "<h2>Log In</h2>")
	assert.NotContains(t, r, "Logged in as")
	assert.NotContains(t, r, "Log out")
}
