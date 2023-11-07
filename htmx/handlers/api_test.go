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
	assert.NotContains(t, r, "Alice")
	assert.Contains(t, r, "Bob")
	assert.Contains(t, r, "Carl")
}

func TestRenderFollowing(t *testing.T) {
	_, e, loginCookie := login()
	req := httptest.NewRequest(http.MethodGet, "/following", strings.NewReader(""))
	req.Header.Set("Cookie", loginCookie)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	r := rec.Body.String()

	assert.Contains(t, r, "Following")
	assert.NotContains(t, r, "Alice")
	assert.Contains(t, r, "Bob")
	assert.NotContains(t, r, "Carl")
}
