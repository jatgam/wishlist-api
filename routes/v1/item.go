package v1

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/jatgam/wishlist-api/metrics"
	"github.com/jatgam/wishlist-api/microservice"
	"github.com/jatgam/wishlist-api/service"
	"github.com/jatgam/wishlist-api/types"
)

type addItemForm struct {
	Name string `form:"name" binding:"required,notblank,alphanumunicode"`
	URL  string `form:"url" binding:"required,notblank,url"`
	Rank int    `form:"rank" binding:"required,notblank,numeric"`
}

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

func isAuthorized(c *gin.Context, requiredlevel float64) bool {
	mylogger := microservice.GetLogger(c)
	claims := jwt.ExtractClaims(c)
	usrlvl, usrlvlFound := claims["userlevel"]
	if !usrlvlFound {
		mylogger.Debugf("User Level Not found in jwt claims: %v", claims)
		return false
	}
	if usrlvlInt, ok := usrlvl.(float64); ok {
		if usrlvlInt == requiredlevel {
			return true
		} else {
			mylogger.Debugf("User Level Didnt Match: %v, Required: %v", usrlvl, requiredlevel)
			return false
		}
	}
	mylogger.Debugf("Failed to Check User Level: %v, Required: %v", usrlvl, requiredlevel)
	return false
}

func addItem(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	var newItem addItemForm
	if !isAuthorized(c, 9) {
		mylogger.Debug("AddItem: Unauthorized")
		types.WriteResponse(c, http.StatusUnauthorized, "Unathorized Access")
		return
	}

	if err := c.ShouldBind(&newItem); err != nil {
		mylogger.Debugf("addItem Failed Form Data Validation")
		// Metric?
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := service.AddItem(strings.TrimSpace(newItem.Name), strings.TrimSpace(newItem.URL), newItem.Rank, mylogger); err != nil {
		mylogger.Error("Item Add Failed.")
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Debug("Authorized to add Items")
	types.WriteResponse(c, http.StatusOK, "Item Created.")
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
