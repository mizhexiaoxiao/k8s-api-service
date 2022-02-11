package routers

import (
	"github.com/gin-gonic/gin"
	istiov1 "github.com/mizhexiaoxiao/k8s-api-service/api/v1/istio"
)

func addIstioRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/istio")

	router.GET("/:cluster/vs", istiov1.GetVirtualServices)
	router.POST("/:cluster/vs", istiov1.PostVirtualService)
	router.GET("/:cluster/vs/:namespace/:vsName", istiov1.GetVirtualService)
	router.PUT("/:cluster/vs/:namespace/:vsName", istiov1.PutVirtualService)
	router.DELETE("/:cluster/vs/:namespace/:vsName", istiov1.DeleteVirtualService)

	router.GET("/:cluster/dr", istiov1.GetDestinationRules)
	router.POST("/:cluster/dr", istiov1.PostDestinationRule)
	router.GET("/:cluster/dr/:namespace/:drName", istiov1.GetDestinationRule)
	router.PUT("/:cluster/dr/:namespace/:drName", istiov1.PutDestinationRule)
	router.DELETE("/:cluster/dr/:namespace/:drName", istiov1.DeleteDestinationRule)
}
