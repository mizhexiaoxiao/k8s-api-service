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
	router.GET("/:cluster/pods/:namespace/:podName/:containerName/download_log", k8sv1.DownloadPodContainerLog)
	router.GET("/:cluster/pods/:namespace/:podName", k8sv1.GetPod)
	router.DELETE("/:cluster/pods/:namespace/:podName", k8sv1.DeletePod)

	router.GET("/:cluster/deployments", k8sv1.GetDeployments)
	router.GET("/:cluster/deployments/:namespace/:deploymentName", k8sv1.GetDeployment)
	router.POST("/:cluster/deployments", k8sv1.PostDeployment)
	router.POST("/:cluster/deployments/:namespace/:deploymentName", k8sv1.DeploymentDoAction)
	router.DELETE("/:cluster/deployments/:namespace/:deploymentName", k8sv1.DeleteDeployment)
	router.PUT("/:cluster/deployments/:namespace/:deploymentName", k8sv1.PutDeployment)
	router.PATCH("/:cluster/deployments/:namespace/:deploymentName", k8sv1.PatchDeployment)
	router.GET("/:cluster/deployment_status/:namespace/:deploymentName", k8sv1.GetDeploymentStatus)
	router.GET("/:cluster/deployment_pods/:namespace/:deploymentName", k8sv1.GetDeploymentPods)

	router.GET("/:cluster/services", k8sv1.GetServices)
	router.GET("/:cluster/services/:namespace/:serviceName", k8sv1.GetService)

	router.GET("/:cluster/jobs", k8sv1.GetJobs)
	router.GET("/:cluster/jobs/:namespace/:jobName", k8sv1.GetJob)
	router.DELETE("/:cluster/jobs/:namespace/:jobName", k8sv1.DeleteJob)

	router.GET("/:cluster/cronjobs", k8sv1.GetCronJobs)
	router.POST("/:cluster/cronjobs", k8sv1.PostCronJob)
	router.GET("/:cluster/cronjobs/:namespace/:cronjobName", k8sv1.GetCronJob)
	router.PUT("/:cluster/cronjobs/:namespace/:cronjobName", k8sv1.PutCronJob)
	router.DELETE("/:cluster/cronjobs/:namespace/:cronjobName", k8sv1.DeleteCronJob)

	router.GET("/:cluster/events", k8sv1.GetEvents)

	router.GET("/:cluster/nodes", k8sv1.GetNodes)
	router.GET("/:cluster/namespaces", k8sv1.GetNamespaces)

	router.POST("/:cluster/horizontalpodautoscalers", k8sv1.PostHorizontalPodAutoScaler)
	router.GET("/:cluster/horizontalpodautoscalers", k8sv1.GetHorizontalPodAutoScalerList)
	router.GET("/:cluster/horizontalpodautoscalers/:namespace/:name", k8sv1.GetHorizontalPodAutoScaler)
	router.PUT("/:cluster/horizontalpodautoscalers/:namespace/:name", k8sv1.PutHorizontalPodAutoScaler)
	router.DELETE("/:cluster/horizontalpodautoscalers/:namespace/:name", k8sv1.DeleteHorizontalPodAutoScaler)

	router.POST("/:cluster/configmaps", k8sv1.PostConfigmap)
	router.GET("/:cluster/configmaps", k8sv1.GetConfigmapList)
	router.GET("/:cluster/configmaps/:namespace/:name", k8sv1.GetConfigmap)
	router.PUT("/:cluster/configmaps/:namespace/:name", k8sv1.PutConfigmap)
	router.DELETE("/:cluster/configmaps/:namespace/:name", k8sv1.DeleteConfigmap)

	router.GET("/:cluster/crd/:group/:version/:resource/:namespace/:name", k8sv1.GetCRD)
	router.POST("/:cluster/crd/:group/:version/:resource", k8sv1.PostCRD)
	router.PUT("/:cluster/crd/:group/:version/:resource/:namespace/:name", k8sv1.PutCRD)
	router.DELETE("/:cluster/crd/:group/:version/:resource/:namespace/:name", k8sv1.DeleteCRD)
}
