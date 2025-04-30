package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LimitBodySize middleware limits the size of the request body
func LimitBodySize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
