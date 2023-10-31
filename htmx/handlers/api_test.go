package handlers

import (
	"fmt"
	"htmx/db"
	"htmx/router"
	"htmx/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRequest(method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
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

	return writer
}

// TODO: Probably just make the rest of the handlers tests hit the endpoints
// instead of calling the handlers
func TestRenderIndex(t *testing.T) {
	rec := makeRequest(http.MethodGet, "/", "", nil)
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

func TestRenderRegister(t *testing.T) {
	rec := makeRequest(http.MethodGet, "/register", "", nil)
	r := rec.Body.String()
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/register\" hx-swap=\"outerHTML\"")
}

func TestRenderLogin(t *testing.T) {
	rec := makeRequest(http.MethodGet, "/login", "", nil)
	r := rec.Body.String()
	fmt.Println(r)
	assert.Contains(t, r, "html")
	assert.Contains(t, r, "nav")
	assert.Contains(t, r, "form hx-post=\"/login\" hx-swap=\"outerHTML\"")
}
