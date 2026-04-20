package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGzipRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Gzip())
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	r.POST("/echo", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(http.StatusOK, string(body))
	})
	return r
}

func TestGzip_PlainRequest(t *testing.T) {
	r := newGzipRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "hello world", w.Body.String())
}

func TestGzip_AcceptEncodingGzip(t *testing.T) {
	r := newGzipRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(w.Body)
	require.NoError(t, err)
	defer gr.Close()
	body, err := io.ReadAll(gr)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(body))
}

func TestGzip_CompressedRequestBody(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte("compressed input"))
	gz.Close()

	r := newGzipRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/echo", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "compressed input", w.Body.String())
}
