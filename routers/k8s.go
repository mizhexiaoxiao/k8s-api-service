package routers

import (
	"github.com/gin-gonic/gin"
	k8sv1 "github.com/mizhexiaoxiao/k8s-api-service/api/v1/k8s"
)

func addK8sRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/k8s")

	router.GET("/:cluster/pods", k8sv1.GetPods)
	router.GET("/:cluster/pods/:namespace/:podName/ssh", k8sv1.PodWebSSH)
	router.GET("/:cluster/pods/:namespace/:podName/log", k8sv1.GetPodLog)
	router.GET("/:cluster/pods/:namespace/:podName", k8sv1.GetPod)

	router.GET("/:cluster/deployments", k8sv1.GetDeployments)
	router.GET("/:cluster/deployments/:namespace/:deploymentName", k8sv1.GetDeployment)
	router.DELETE("/:cluster/deployments/:namespace/:deploymentName", k8sv1.DeleteDeployment)
	router.PUT("/:cluster/deployments/:namespace/:deploymentName", k8sv1.PutDeployment)
	router.GET("/:cluster/deployment_status/:namespace/:deploymentName", k8sv1.GetDeploymentStatus)
	router.GET("/:cluster/deployment_pods/:namespace/:deploymentName", k8sv1.GetDeploymentPods)

	router.GET("/:cluster/services", k8sv1.GetServices)
	router.GET("/:cluster/services/:namespace/:serviceName", k8sv1.GetService)

	router.GET("/:cluster/events", k8sv1.GetEvents)

	router.GET("/:cluster/nodes", k8sv1.GetNodes)
	router.GET("/:cluster/namespaces", k8sv1.GetNamespaces)
}
