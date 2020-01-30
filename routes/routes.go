package routes

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	v1 "github.com/jatgam/wishlist-api/routes/v1"
)

func health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func SetupRoutes(router *gin.RouterGroup, ginjwt *jwt.GinJWTMiddleware) {
	router.GET("/health", health)
	v1.SetupV1Routes(router, ginjwt) // V1 Technically was never versioned, so all routes are at the root.
}
