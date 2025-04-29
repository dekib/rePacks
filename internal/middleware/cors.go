package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Set default security headers
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// 2. Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 3. Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// 4. Allow specific origins when configured
		if origin := c.GetHeader("Origin"); origin != "" {
			allowedOrigins := getConfiguredOrigins()
			if isOriginAllowed(origin, allowedOrigins) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		}

		c.Next()
	}
}

// Helper functions
func getConfiguredOrigins() []string {
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		return strings.Split(envOrigins, ",")
	}
	return []string{"*"} // Fallback to allow all for development
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, o := range allowedOrigins {
		if o == "*" || strings.EqualFold(o, origin) {
			return true
		}
	}
	return false
}
