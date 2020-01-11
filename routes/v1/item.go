package v1

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/jatgam/wishlist-api/metrics"
	"github.com/jatgam/wishlist-api/microservice"
	"github.com/jatgam/wishlist-api/service"
	"github.com/jatgam/wishlist-api/types"
)

func getWantedItems(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	items, err := service.GetWantedItems(mylogger)

	if err != nil {
		mylogger.Error("Failed to get Wanted Items")
		metrics.ItemErrors.WithLabelValues(metrics.ItemGetError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Info("Got Wanted Items")
	types.WriteItemResponse(c, http.StatusOK, "Got a list of Wanted items", items)
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
