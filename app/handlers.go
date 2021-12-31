package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{Code: http.StatusNotFound, Msg: "No route found", Data: nil})
	return
}

func HandleNotMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, Response{Code: http.StatusMethodNotAllowed, Msg: "Method not allowed", Data: nil})
	return
}
