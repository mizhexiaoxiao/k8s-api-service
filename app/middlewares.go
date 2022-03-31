package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/config"
)

//简单校验
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := Gin{C: c}
		url := appG.C.Request.RequestURI
		if url == "/api/v1/health_check" || strings.HasPrefix(url, "/swagger/") {
			c.Next()
			return
		}
		appKey := appG.C.GetHeader("appKey")
		appSecret := appG.C.GetHeader("appSecret")
		configSecret := config.GetString("caller" + "." + appKey + ".secret")
		if configSecret == "" {
			appG.Fail(http.StatusInternalServerError, errors.New(fmt.Sprintf("Please carry appKey and appSecret in the request header")), nil)
			log.Println(errors.New(fmt.Sprintf("Please carry appKey and appSecret in the request header")))
			c.Abort()
			return
		}
		if appSecret != configSecret {
			appG.Fail(http.StatusUnauthorized, errors.New("Authentication failed"), nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
