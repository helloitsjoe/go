package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderIndex(t *testing.T) {
	rec, _ := makeRequest(http.MethodGet, "/", "", nil)
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

// func TestRender404(t *testing.T) {
// 	rec, _ := makeRequest(http.MethodGet, "/nope", "", nil)
// 	r := rec.Body.String()
// 	assert.Contains(t, r, "html")
// 	assert.Contains(t, r, "nav")
// 	assert.Contains(t, r, "404")
// }

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
	_, e, loginCookie := login()
	req := httptest.NewRequest(http.MethodGet, "/followers", strings.NewReader(""))
	req.Header.Set("Cookie", loginCookie)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	r := rec.Body.String()

	assert.Contains(t, r, "Followers")
	assert.NotContains(t, r, "alice")
	assert.Contains(t, r, "bob")
	assert.Contains(t, r, "carl")
	assert.Contains(t, r, `href="/user/bob"`)
	assert.Contains(t, r, `href="/user/carl"`)
}

func TestRenderFollowing(t *testing.T) {
	_, e, loginCookie := login()
	req := httptest.NewRequest(http.MethodGet, "/following", strings.NewReader(""))
	req.Header.Set("Cookie", loginCookie)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	r := rec.Body.String()

	assert.Contains(t, r, "Following")
	assert.NotContains(t, r, "alice")
	assert.Contains(t, r, "bob")
	assert.Contains(t, r, `href="/user/bob"`)
	assert.NotContains(t, r, "carl")
}

// TODO: Test button: follow if not following, unfollow otherwise
func TestRenderUser(t *testing.T) {
	_, e, loginCookie := login()
	req := httptest.NewRequest(http.MethodGet, "/user/bob", strings.NewReader(""))
	req.Header.Set("Cookie", loginCookie)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	r := rec.Body.String()

	assert.Contains(t, r, "bob's books")
	assert.Contains(t, r, "checked out")
	assert.Contains(t, r, "available")
	assert.Contains(t, r, "Logged in as alice")
}
