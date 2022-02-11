package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
)

func HealthCheck(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Success(http.StatusOK, "ok", nil)
	return
}
