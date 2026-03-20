package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
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
func TestHellohandler_Metrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/api/hello", hellohandler)

	// reset global metrics to avoid cross-test contamination
	httpRequestsTotal.Reset()
	httpRequestDuration.Reset()

	// 1st request: error path
	reqErr, err := http.NewRequest(http.MethodGet, "/api/hello?name=err", nil)
	assert.NoError(t, err)
	wErr := httptest.NewRecorder()
	r.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)

	// 2nd request: success path
	reqOK, err := http.NewRequest(http.MethodGet, "/api/hello?name=world", nil)
	assert.NoError(t, err)
	wOK := httptest.NewRecorder()
	r.ServeHTTP(wOK, reqOK)
	assert.Equal(t, http.StatusOK, wOK.Code)

	// assert counters
	assert.Equal(t, float64(1), testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/api/hello", "500")))
	assert.Equal(t, float64(1), testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/api/hello", "200")))

	// assert histogram has recorded 2 observations (one per response path)
	count := testutil.CollectAndCount(httpRequestDuration)
	assert.Equal(t, 2, count)

}
