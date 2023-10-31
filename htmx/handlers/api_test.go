package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
