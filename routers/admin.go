package routers

import (
	"github.com/gin-gonic/gin"
	adminv1 "github.com/mizhexiaoxiao/k8s-api-service/api/v1/admin"
)

func addAdminRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/admin")

	router.GET("/clusters", adminv1.ListCluster)
	router.POST("/clusters", adminv1.PostCluster)
	router.PUT("/clusters/:id", adminv1.PutCluster)
	router.GET("/clusters/:id", adminv1.GetCluster)
	router.DELETE("/clusters/:id", adminv1.DeleteCluster)
	router.POST("/testConnectclusters/", adminv1.TestConnectCluster)
}
