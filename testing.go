package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TesthelloHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	test := []struct {
		name    string
		path    string
		code    int
		contain string
	}{
		{name: "hello world", path: "/hello", code: 200, contain: "Hello, World!"},
		{name: "not found", path: "/notfound", code: 404, contain: "404 page not found"},
		{name: "invalid method", path: "/hello", code: 405, contain: "Method Not Allowed"},
	}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = request
			hellohandler(c)

			assert.Equal(t, tc.code, w.Code)
			assert.Contains(t, w.Body.String(), tc.contain)
		})
	}
}
