package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Log("echo")
}

func TestHellohandler(t *testing.T) {
	//TODO
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/api/hello")

	test := []struct {
		name string
		code int
		want string
	}{
		{name: "world", code: 200, want: "hello world"},
		{name: "err", code: 500, want: "internal server error"},
		{name: "", code: 200, want: "hello world"},
	}
	for _, tt := range test {

		t.Logf("test with name: %s", tt.name)

		req, err := http.NewRequest(http.MethodGet, "/api/hello?name="+tt.name, nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
