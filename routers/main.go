package routers

import (
	"github.com/gin-gonic/gin"
	adminv1 "github.com/mizhexiaoxiao/k8s-api-service/api/v1/admin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	docs "github.com/mizhexiaoxiao/k8s-api-service/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.HandleMethodNotAllowed = true
	r.NoMethod(app.HandleNotMethod)
	r.NoRoute(app.HandleNotFound)
	//Authentication
	r.Use(app.Auth())
	// swagger config
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/health_check", adminv1.HealthCheck)
		// k8s api
		addK8sRoutes(apiv1)

		// istio api
		addIstioRoutes(apiv1)

		// cluster admin
		addAdminRoutes(apiv1)

	}
	return r
}
