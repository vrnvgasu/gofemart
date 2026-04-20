package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const gzipHeader = "gzip"

type compressWriter struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w gin.ResponseWriter) *compressWriter {
	return &compressWriter{ResponseWriter: w, zw: gzip.NewWriter(w)}
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.ResponseWriter.Header().Set("Content-Encoding", gzipHeader)
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{r: r, zr: zr}, nil
}

func (c *compressReader) Read(p []byte) (int, error) { return c.zr.Read(p) }
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Gzip возвращает middleware для поддержки gzip-сжатия.
// Сжимает ответ, если клиент поддерживает gzip (Accept-Encoding: gzip).
// Распаковывает тело запроса, если оно сжато (Content-Encoding: gzip).
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), gzipHeader) {
			cw := newCompressWriter(c.Writer)
			defer cw.Close()
			c.Writer = cw
		}

		if strings.Contains(c.Request.Header.Get("Content-Encoding"), gzipHeader) {
			cr, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer cr.Close()
			c.Request.Body = cr
		}

		c.Next()
	}
}
