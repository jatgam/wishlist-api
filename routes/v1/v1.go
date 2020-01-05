package v1

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// SetupV1Routes sets up the routes for the V1 API
func SetupV1Routes(router *gin.RouterGroup, ginjwt *jwt.GinJWTMiddleware) {
	userGroup := router.Group("/user")
	setupUserRoutes(userGroup, ginjwt)

	itemGroup := router.Group("/item")
	setupItemRoutes(itemGroup, ginjwt)
}
