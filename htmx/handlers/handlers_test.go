package handlers

import (
	"encoding/json"
	"fmt"
	"htmx/db"
	"htmx/router"
	"htmx/types"
	"htmx/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var headers = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

// TODO: makeInitialRequest/makeFollowUpRequest

func makeRequest(method, path, body string, headers map[string]string) (*httptest.ResponseRecorder, *echo.Echo) {
	e := router.New("../")
	d := db.CreateDB()
	user.SeedUsers(d)
	Register(e, d)
	req := httptest.NewRequest(method, path, strings.NewReader(body))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	writer := httptest.NewRecorder()
	e.ServeHTTP(writer, req)

	return writer, e
}

func login() (*httptest.ResponseRecorder, *echo.Echo, string) {
	rec, e := makeRequest(http.MethodPost, "/login", "username=Alice&password=bar", headers)

	c := rec.Header().Get("Set-Cookie")
	cookie := strings.Split(c, ";")[0]
	cookieVal := strings.Split(cookie, "=")[1]
	loginCookie := fmt.Sprintf("uuid=%s", cookieVal)

	return rec, e, loginCookie
}

func TestGetUsersHtmx(t *testing.T) {
	rec, _ := makeRequest(http.MethodGet, "/users", "", nil)
	r := rec.Body.String()
	assert.Contains(t, r, "alice")
	assert.Contains(t, r, "bob")
	assert.Contains(t, r, "carl")
}

func TestGetUsersJson(t *testing.T) {
	rec, _ := makeRequest(http.MethodGet, "/users?format=json", "", nil)
	r := rec.Body.Bytes()
	users := []types.User{}
	err := json.Unmarshal(r, &users)
	assert.Nil(t, err)

	alice, bob, carl := users[0], users[1], users[2]

	assert.Equal(t, alice.Username, "alice")
	assert.Equal(t, bob.Username, "bob")
	assert.Equal(t, carl.Username, "carl")

	// Bit annoying way to check that the Ids are unique
	assert.NotEqual(t, alice.UUID, bob.UUID)
	assert.NotEqual(t, bob.UUID, carl.UUID)
	assert.NotEqual(t, alice.UUID, carl.UUID)
}

func TestRegisterUserNoInputFail(t *testing.T) {
	rec, _ := makeRequest(http.MethodPost, "/register", "", nil)
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Name and password must be provided")
}

func TestRegisterUserExistsFail(t *testing.T) {
	rec, _ := makeRequest(http.MethodPost, "/register", "username=alice&password=bar", headers)
	r := rec.Body.String()

	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "That username is already taken")

	cookie := rec.Header().Get("Set-Cookie")
	assert.NotContains(t, cookie, `uuid=`)

	assert.NotContains(t, r, "Hello there alice")
	assert.NotContains(t, r, "Log out")
}

func TestRegisterUserValid(t *testing.T) {
	rec, _ := makeRequest(http.MethodPost, "/register", "username=NewUser&password=123", headers)
	r := rec.Body.String()
	cookie := rec.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, `uuid=`)
	assert.Contains(t, r, "<span id=\"num-users\" hx-swap-oob=\"true\">4</span>")
	assert.Contains(t, r, `
<span id="logged-in-as" hx-swap-oob="true">
  <span>Logged in as newuser</span>
  <button hx-post="/logout" class="text-btn">Log out</button>
</span>`)
	assert.Contains(t, r, "Hello there newuser")
	assert.Contains(t, r, "Log out")
}

func TestLogin(t *testing.T) {
	rec, _ := makeRequest(echo.POST, "/login", "username=Alice&password=bar", headers)
	r := rec.Body.String()
	cookie := rec.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, "uuid=")

	assert.Contains(t, r, "<span id=\"num-users\" hx-swap-oob=\"true\">3</span>")
	assert.Contains(t, r, `
<span id="logged-in-as" hx-swap-oob="true">
  <span>Logged in as alice</span>
  <button hx-post="/logout" class="text-btn">Log out</button>
</span>`)
	// Clear error on a successful login
	assert.Contains(t, r, "<div id=\"error\" hx-swap-oob=\"true\"></div>")
	assert.Contains(t, r, "Hello there alice")
	assert.Contains(t, r, "Log out")
}

// TODO: convert some of these to user unit tests
func TestLoginError(t *testing.T) {
	rec, _ := makeRequest(echo.POST, "/login", "username=Alice&password=wrong-password", headers)
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Incorrect login credentials")
}

func TestLoginNoUsernamePassword(t *testing.T) {
	rec, _ := makeRequest(echo.POST, "/login", "username=&password=", headers)
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Name and password must be provided")
}

func TestLoginNoAccount(t *testing.T) {
	rec, _ := makeRequest(echo.POST, "/login", "username=Not a user&password=bar", headers)
	r := rec.Body.String()
	assert.Equal(t, rec.Header().Get("HX-Retarget"), "#error")
	assert.Equal(t, rec.Header().Get("HX-Reswap"), "innerHTML")
	assert.Contains(t, r, "div class=\"error\"")
	assert.Contains(t, r, "Incorrect login credentials")
}

func TestLogout(t *testing.T) {
	rec, _ := makeRequest(echo.POST, "/logout", "", nil)
	r := rec.Body.String()
	assert.Contains(t, rec.Header().Get("Set-Cookie"), "uuid=; Max-Age=0; HttpOnly")
	assert.Contains(t, r, "<h2>Log In</h2>")
	assert.NotContains(t, r, "Hello there")
	assert.NotContains(t, r, "Log out")
}
