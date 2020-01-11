package types

import (
	"github.com/gin-gonic/gin"
)

type (
	GenericResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	Items struct {
		Name string `json:"name"`
		Url  string `json:"url"`
		Rank int    `json:"rank"`
	}
	GetItemsResponse struct {
		GenericResponse
		Items *[]Items `json:"items"`
	}
)

// WriteResponse will create the generic json response, and set the gin
// response status, and end the request to make sure no further handlers
// are called.
func WriteResponse(c *gin.Context, code int, message string) {
	resp := GenericResponse{Code: code, Message: message}
	c.JSON(code, resp)
	c.Abort()
}

func WriteItemResponse(c *gin.Context, code int, message string, items *[]Items) {
	resp := GetItemsResponse{GenericResponse{Code: code, Message: message}, items}
	c.JSON(code, resp)
	c.Abort()
}
