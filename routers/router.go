package routers

import (
	"github.com/gin-gonic/gin"
	adminv1 "github.com/mizhexiaoxiao/k8s-api-service/api/v1/admin"
	k8sv1 "github.com/mizhexiaoxiao/k8s-api-service/api/v1/k8s"
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
		// k8s api
		apiv1.GET("/health_check", k8sv1.HealthCheck)
		apiv1.GET("/k8s/:cluster/pods", k8sv1.GetPods)
		apiv1.GET("/k8s/:cluster/pods/:namespace/:podName/log", k8sv1.GetPodLog)
		apiv1.GET("/k8s/:cluster/pods/:namespace/:podName", k8sv1.GetPod)

		apiv1.GET("/k8s/:cluster/deployments", k8sv1.GetDeployments)
		apiv1.GET("/k8s/:cluster/deployments/:namespace/:deploymentName", k8sv1.GetDeployment)
		apiv1.PUT("/k8s/:cluster/deployments/:namespace/:deploymentName", k8sv1.PutDeployment)
		apiv1.GET("/k8s/:cluster/deployment_status/:namespace/:deploymentName", k8sv1.GetDeploymentStatus)
		apiv1.GET("/k8s/:cluster/deployment_pods/:namespace/:deploymentName", k8sv1.GetDeploymentPods)

		apiv1.GET("/k8s/:cluster/services", k8sv1.GetServices)
		apiv1.GET("/k8s/:cluster/services/:namespace/:serviceName", k8sv1.GetService)

		apiv1.GET("/k8s/:cluster/events", k8sv1.GetEvents)

		apiv1.GET("/k8s/:cluster/nodes", k8sv1.GetNodes)
		apiv1.GET("/k8s/:cluster/namespaces", k8sv1.GetNamespaces)

		// cluster admin
		apiv1.GET("/admin/clusters", adminv1.ListCluster)
		apiv1.POST("/admin/clusters", adminv1.PostCluster)
		apiv1.PUT("/admin/clusters/:id", adminv1.PutCluster)
		apiv1.GET("/admin/clusters/:id", adminv1.GetCluster)
		apiv1.DELETE("/admin/clusters/:id", adminv1.DeleteCluster)
		apiv1.POST("/admin/testConnectclusters/", adminv1.TestConnectCluster)

	}
	return r
}
