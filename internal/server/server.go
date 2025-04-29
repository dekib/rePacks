package server

import (
	//csrf "github.com/utrack/gin-csrf"
	"os"
	"time"

	v1 "github.com/dekib/rePacks/controllers/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	r := gin.Default()

	// Load HTML from the ui directory
	r.LoadHTMLGlob("ui/*.html")

	// CORS configuration
	if corsEnv := os.Getenv("CORS"); corsEnv != "" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{corsEnv}
		r.Use(cors.New(config))
	}

	//// Add the CSRF middleware
	//r.Use(csrf.Middleware(csrf.Options{
	//	Secret: "a-32-byte-long-secret-key!",
	//	ErrorFunc: func(c *gin.Context) {
	//		c.String(400, "CSRF token mismatch")
	//		c.Abort()
	//	},
	//}))

	// Global timeout middleware
	r.Use(timeoutMiddleware())

	// Routes
	v1.Routes(r)

	// Serve HTML UI at root
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	// Health check
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(8*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	c.JSON(504, gin.H{
		"error": "Request timeout",
	})
}
