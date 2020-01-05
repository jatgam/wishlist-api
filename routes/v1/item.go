package v1

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	// "github.com/jatgam/wishlist-api/microservice"
	// "github.com/jatgam/wishlist-api/models"
)

func getWantedItems(c *gin.Context) {

}

func addItem(c *gin.Context) {

}

func getAllItems(c *gin.Context) {

}

func getReservedItems(c *gin.Context) {

}

func deleteItem(c *gin.Context) {

}

func reserveItem(c *gin.Context) {

}

func unReserveItem(c *gin.Context) {

}

func editItemRank(c *gin.Context) {

}

func setupItemRoutes(router *gin.RouterGroup, ginjwt *jwt.GinJWTMiddleware) {
	authMiddleware := ginjwt.MiddlewareFunc()
	router.GET("", getWantedItems)
	router.POST("", authMiddleware, addItem)
	router.GET("/all", authMiddleware, getAllItems)
	router.GET("/reserved", authMiddleware, getReservedItems)
	router.DELETE("/id/:itemID", authMiddleware, deleteItem)
	router.POST("/id/:itemID/reserve", authMiddleware, reserveItem)
	router.POST("/id/:itemID/unreserve", authMiddleware, unReserveItem)
	router.POST("/id/:itemID/rank/:rank", authMiddleware, editItemRank)
}
