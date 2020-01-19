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

type itemRankEditURI struct {
	ItemID int `uri:"itemID" binding:"required,numeric,notblank"`
	Rank   int `uri:"rank" binding:"required,numeric,notblank"`
}

type itemURI struct {
	ItemID int `uri:"itemID" binding:"required,numeric,notblank"`
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

func getAuthenticatedUsersID(c *gin.Context) (int, error) {
	mylogger := microservice.GetLogger(c)
	claims := jwt.ExtractClaims(c)
	userID, userIDFound := claims["id"]
	if !userIDFound {
		mylogger.Debugf("User ID Not found in jwt claims: %v", claims)
		return 0, types.ErrDeterminingUserIDFromJWT
	}
	if userIDF, ok := userID.(float64); ok {
		userIDInt := int(userIDF)
		mylogger.Debugf("User IDF: %v, IDI: %v", userIDF, userIDInt)
		return userIDInt, nil
	}
	mylogger.Errorf("User ID in jwt is not a number: %v", userID)
	return 0, types.ErrDeterminingUserIDFromJWT
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
		metrics.ItemErrors.WithLabelValues(metrics.ItemAddDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := service.AddItem(strings.TrimSpace(newItem.Name), strings.TrimSpace(newItem.URL), newItem.Rank, mylogger); err != nil {
		mylogger.Error("Item Add Failed.")
		metrics.ItemErrors.WithLabelValues(metrics.ItemAddError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Debug("Authorized to add Items")
	types.WriteResponse(c, http.StatusOK, "Item Created.")
}

func getAllItems(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	if !isAuthorized(c, 9) {
		mylogger.Debug("GetAllItems: Unauthorized")
		types.WriteResponse(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}

	items, err := service.GetAllItems(mylogger)

	if err != nil {
		mylogger.Error("Failed to get All Items")
		metrics.ItemErrors.WithLabelValues(metrics.ItemGetError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Info("Got All Items")
	types.WriteItemResponse(c, http.StatusOK, "Got a list of items", items)

}

func getReservedItems(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	userID, err := getAuthenticatedUsersID(c)
	if err != nil {
		mylogger.Errorf("Failed to get reserved items: %s", err.Error())
		metrics.ItemErrors.WithLabelValues(metrics.ItemGetError).Inc()
		types.WriteResponse(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}
	items, err := service.GetReservedItems(userID, mylogger)

	if err != nil {
		mylogger.Error("Failed to get Reserved Items")
		metrics.ItemErrors.WithLabelValues(metrics.ItemGetError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Info("Got Reserved Items")
	types.WriteItemResponse(c, http.StatusOK, "Got a list of reserved items", items)
}

func deleteItem(c *gin.Context) {
	mylogger := microservice.GetLogger(c)

	if !isAuthorized(c, 9) {
		mylogger.Debug("DeleteItem: Unauthorized")
		types.WriteResponse(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}

	var itemInfo itemURI
	if err := c.ShouldBindUri(&itemInfo); err != nil {
		mylogger.Debug("Delete Item Data Validation Error")
		metrics.ItemErrors.WithLabelValues(metrics.ItemDeleteValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if deleteError := service.DeleteItem(itemInfo.ItemID, mylogger); deleteError != nil {
		mylogger.Error("Failed to Delete Item")
		metrics.ItemErrors.WithLabelValues(metrics.ItemDeleteError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, deleteError.Error())
		return
	}

	mylogger.Info("Item Deleted")
	types.WriteResponse(c, http.StatusOK, "Item Deleted")
}

func reserveItem(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	userID, err := getAuthenticatedUsersID(c)
	if err != nil {
		mylogger.Errorf("Failed to ReserveItem: %s", err.Error())
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditError).Inc()
		types.WriteResponse(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}
	var itemInfo itemURI
	if err := c.ShouldBindUri(&itemInfo); err != nil {
		mylogger.Debug("Reserve Item Data Validation Error")
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	reserveErr := service.ReserveItem(userID, itemInfo.ItemID, mylogger)
	if reserveErr != nil {
		mylogger.Error("Failed to Reserve Item")
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Info("Item Reserved")
	types.WriteResponse(c, http.StatusOK, "Item Reserved")
}

func unReserveItem(c *gin.Context) {
	mylogger := microservice.GetLogger(c)
	userID, err := getAuthenticatedUsersID(c)
	if err != nil {
		mylogger.Errorf("Failed to UnReserveItem: %s", err.Error())
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditError).Inc()
		types.WriteResponse(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}
	var itemInfo itemURI
	if err := c.ShouldBindUri(&itemInfo); err != nil {
		mylogger.Debug("UnReserveItem Data Validation Error")
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	unreserveErr := service.UnReserveItem(userID, itemInfo.ItemID, mylogger)
	if unreserveErr != nil {
		mylogger.Error("Failed to UnReserve Item")
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Info("Item UnReserved")
	types.WriteResponse(c, http.StatusOK, "Item UnReserved")
}

func editItemRank(c *gin.Context) {
	mylogger := microservice.GetLogger(c)

	if !isAuthorized(c, 9) {
		mylogger.Debug("EditItemRank: Unauthorized")
		types.WriteResponse(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}

	var itemInfo itemRankEditURI
	if err := c.ShouldBindUri(&itemInfo); err != nil {
		mylogger.Debug("Edit Item Rank Data Validation Error")
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditDataValidationError).Inc()
		types.WriteResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	err := service.EditItemRank(itemInfo.ItemID, itemInfo.Rank, mylogger)

	if err != nil {
		mylogger.Error("Failed to Edit Item")
		metrics.ItemErrors.WithLabelValues(metrics.ItemEditError).Inc()
		types.WriteResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	mylogger.Info("Item Rank Edited")
	types.WriteResponse(c, http.StatusOK, "Item Rank Updated")
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
