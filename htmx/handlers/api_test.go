package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRenderIndex(t *testing.T) {
	rec, _ := makeRequest(http.MethodGet, "/", "", nil)
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

func TestRenderRegister(t *testing.T) {
	rec, _ := makeRequest(http.MethodGet, "/register", "", nil)
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

func TestRenderLogin(t *testing.T) {
	rec, _ := makeRequest(http.MethodGet, "/login", "", nil)
	r := rec.Body.String()
	fmt.Println(r)
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/login\" hx-swap=\"outerHTML\"")
}

func TestRenderFollowers(t *testing.T) {
	loginRec, e := makeRequest(echo.POST, "/login", "username=Alice&password=bar", headers)
	c := loginRec.Header().Get("Set-Cookie")
	cookie := strings.Split(c, ";")[0]
	cookieVal := strings.Split(cookie, "=")[1]

	assert.NotZero(t, cookieVal)

	req := httptest.NewRequest(http.MethodGet, "/followers", strings.NewReader(""))
	cookieHeader := fmt.Sprintf("uuid=%s", cookieVal)
	req.Header.Set("Cookie", cookieHeader)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	r := rec.Body.String()

	assert.Contains(t, r, "Followers")
	assert.NotContains(t, r, "Alice")
	assert.Contains(t, r, "Bob")
	assert.Contains(t, r, "Carl")
}

func TestRenderFollowing(t *testing.T) {
	loginRec, e := makeRequest(echo.POST, "/login", "username=Alice&password=bar", headers)
	c := loginRec.Header().Get("Set-Cookie")
	cookie := strings.Split(c, ";")[0]
	cookieVal := strings.Split(cookie, "=")[1]

	assert.NotZero(t, cookieVal)

	req := httptest.NewRequest(http.MethodGet, "/following", strings.NewReader(""))
	cookieHeader := fmt.Sprintf("uuid=%s", cookieVal)
	req.Header.Set("Cookie", cookieHeader)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	r := rec.Body.String()

	assert.Contains(t, r, "Following")
	assert.NotContains(t, r, "Alice")
	assert.Contains(t, r, "Bob")
	assert.NotContains(t, r, "Carl")
}
