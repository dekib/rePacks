package v1

import (
	"github.com/dekib/rePacks/controllers/v1/pack"
	"github.com/dekib/rePacks/internal/middleware"
	"github.com/dekib/rePacks/internal/services"
	"github.com/gin-gonic/gin"
	"time"
)

func Routes(route *gin.Engine) {
	packService := &services.PackService{}

	// Initialize controllers with dependency injection
	packController := pack.NewPackController(packService)

	v1 := route.Group("/api/v1")
	v1.Use(
		middleware.CorsMiddleware(),
		middleware.TimeoutMiddleware(5*time.Second),
	)
	unauthenticatedRoutes := v1.Group("/")
	unauthenticatedRoutes.Use(middleware.ErrorHandler())

	// Setup pack routes
	packRoutes := unauthenticatedRoutes.Group("/packs")
	{
		packRoutes.GET("", packController.GetPacks)              // GET /api/v1/packs?items_no=250
		packRoutes.PUT("/sizes", packController.UpdatePackSizes) // PUT /api/v1/packs/sizes
	}
}
