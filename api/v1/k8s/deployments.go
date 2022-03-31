package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentsQuery struct {
	Namespace string `form:"namespace"`
	Label     string `form:"label"`
}

type DeploymentsUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type DeploymentQuery struct {
	Label string `json:"label" form:"label"`
}

type DeploymentUri struct {
	Cluster        string `uri:"cluster" binding:"required"`
	Namespace      string `uri:"namespace" binding:"required"`
	DeploymentName string `uri:"deploymentName" binding:"required"`
}

type DeploymentBody struct {
	Image    string `json:"image" form:"image"`
	Label    string `json:"label" form:"label"`
	Replicas string `json:"replicas" form:"replicas"`
}

// @Summary 查看deployment列表
// @Produce  json
// @Param cluster path string true "Cluster"
// @Param namespace query string true "Namespace"
// @Param label query string false "Label"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/deployments [get]
func GetDeployments(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		q        DeploymentsQuery
		u        DeploymentsUri
		listOpts metav1.ListOptions
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if q.Label == "" {
		listOpts = metav1.ListOptions{}
	} else {
		listOpts = metav1.ListOptions{LabelSelector: q.Label}
	}
	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	deployments, err := k8sClient.ClientV1.AppsV1().Deployments(q.Namespace).List(context.TODO(), listOpts)
	for i := 0; i < len(deployments.Items); i++ {
		deployments.Items[i].CreationTimestamp = metav1.NewTime(deployments.Items[i].CreationTimestamp.Add(8 * time.Hour))
	}

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", deployments)
}

// @Summary 查看deployment
// @accept application/json
// @Produce  application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "DeploymentName"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/deployments/{namespace}/{deploymentName} [get]
func GetDeployment(c *gin.Context) {
	appG := app.Gin{C: c}

	var u DeploymentUri

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	deployment, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Get(context.TODO(), u.DeploymentName, metav1.GetOptions{})
	deployment.CreationTimestamp = metav1.NewTime(deployment.CreationTimestamp.Add(8 * time.Hour))

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", deployment)
}

// PostDeployment
// @Summary 创建deployment
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "DeploymentName"
// @Param deployment body metadata.DeploymentCreate true "Deployment"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/deployments/{namespace}/{deploymentName} [post]
func PostDeployment(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u          DeploymentUri
		deployment appsv1.Deployment
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBind(&deployment); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	deployment.ObjectMeta.Name = u.DeploymentName // 名称赋值

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	operation := k8s.NewDeploymentOperation(k8sClient.ClientV1)
	result, err := operation.Create(u.Namespace, &deployment)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// @Summary 更新deployment
// @Produce  json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "DeploymentName"
// @Param RequestBody body v1.DeploymentBody true "RequestBody"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/deployments/{namespace}/{deploymentName} [put]
func PutDeployment(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		b        DeploymentBody
		u        DeploymentUri
		listOpts metav1.ListOptions
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	if b.Label == "" {
		deployment, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Get(context.TODO(), u.DeploymentName, metav1.GetOptions{})

		// update replicas
		if b.Replicas != "" {
			replicas, err := strconv.ParseInt(b.Replicas, 10, 32)
			if err != nil {
				appG.Fail(http.StatusInternalServerError, err, nil)
				return
			}
			r := int32(replicas)
			sc, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).GetScale(context.TODO(), u.DeploymentName, metav1.GetOptions{})
			sc.Spec.Replicas = r
			_, err = k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).UpdateScale(context.TODO(), u.DeploymentName, sc, metav1.UpdateOptions{})
			if err != nil {
				appG.Fail(http.StatusInternalServerError, err, nil)
				return
			}
			appG.Success(http.StatusOK, "deployment replicas update to "+b.Replicas, nil)
			return
		}

		// update image
		if b.Image != "" {
			deployment.Spec.Template.Spec.Containers[0].Image = b.Image
		}

		// force update
		ForceUpdate(deployment)

		_, err = k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
	} else {
		listOpts = metav1.ListOptions{LabelSelector: b.Label}
		deployments, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).List(context.TODO(), listOpts)
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
		for _, deployment := range deployments.Items {
			deployment.Spec.Template.Spec.Containers[0].Image = b.Image
			// force update
			ForceUpdate(&deployment)
			_, err = k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
			if err != nil {
				appG.Fail(http.StatusInternalServerError, err, nil)
				return
			}
		}
	}

	appG.Success(http.StatusOK, "ok", nil)

}

func DeleteDeployment(c *gin.Context) {
	appG := app.Gin{C: c}

	var u DeploymentUri

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	err = k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Delete(context.TODO(), u.DeploymentName, metav1.DeleteOptions{})

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", nil)
}

func GetDeploymentStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u        DeploymentUri
		q        DeploymentQuery
		listOpts metav1.ListOptions
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	if q.Label == "" {
		deployment, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Get(context.TODO(), u.DeploymentName, metav1.GetOptions{})
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
		status, reasons, err := getDeploymentStatus(k8sClient.ClientV1, deployment)
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, reasons)
			return
		}
		if status == 200 {
			appG.Success(http.StatusOK, "ok", nil)
			return
		}
		if status == 308 {
			appG.Fail(http.StatusPermanentRedirect, errors.New("retry"), reasons)
			return
		}
	} else {
		listOpts = metav1.ListOptions{LabelSelector: q.Label}
		deployments, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).List(context.TODO(), listOpts)
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
		if len(deployments.Items) == 0 {
			appG.Fail(http.StatusNotFound, errors.New("deployments not found"), nil)
			return
		}
		for _, deployment := range deployments.Items {
			status, reasons, err := getDeploymentStatus(k8sClient.ClientV1, &deployment)
			if err != nil {
				appG.Fail(status, err, nil)
				return
			}
			if status == 308 {
				appG.Fail(http.StatusPermanentRedirect, errors.New("retry"), reasons)
				return
			}
		}
		appG.Success(http.StatusOK, "ok", nil)
		return
	}

}

func getDeploymentStatus(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) (status int, reasons []string, err error) {
	// 获取pod的状态
	labelSelector := ""
	for key, value := range deployment.Spec.Selector.MatchLabels {
		labelSelector = labelSelector + key + "=" + value + ","
	}
	labelSelector = strings.TrimRight(labelSelector, ",")
	podList, err := clientset.CoreV1().Pods(deployment.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})

	if err != nil {
		return 500, []string{"get pods status error"}, err
	}
	if len(podList.Items) == 0 {
		return 404, []string{"pods not found"}, errors.New("pods not found")
	}
	readyPod := 0
	unavailablePod := 0
	waitingReasons := []string{}
	for _, pod := range podList.Items {
		// 记录等待原因
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil {
				reason := "namespace: " + pod.Namespace + ", pod: " + pod.Name + ", container: " + containerStatus.Name + ", waiting reason: " + containerStatus.State.Waiting.Reason
				waitingReasons = append(waitingReasons, reason)
			}
		}

		podScheduledCondition := GetPodCondition(pod.Status, corev1.PodScheduled)
		initializedCondition := GetPodCondition(pod.Status, corev1.PodInitialized)
		readyCondition := GetPodCondition(pod.Status, corev1.PodReady)
		containersReadyCondition := GetPodCondition(pod.Status, corev1.ContainersReady)

		if pod.Status.Phase == "Running" &&
			podScheduledCondition.Status == "True" &&
			initializedCondition.Status == "True" &&
			readyCondition.Status == "True" &&
			containersReadyCondition.Status == "True" {
			readyPod++
		} else {
			unavailablePod++
		}
	}

	// 根据container状态判定
	if len(waitingReasons) != 0 {
		return 308, waitingReasons, nil
	}

	// 根据pod状态判定
	if int32(readyPod) < *(deployment.Spec.Replicas) ||
		int32(unavailablePod) != 0 {
		return 308, []string{"pods not ready!"}, nil
	}

	// deployment进行状态判定
	availableCondition := GetDeploymentCondition(deployment.Status, appsv1.DeploymentAvailable)
	progressingCondition := GetDeploymentCondition(deployment.Status, appsv1.DeploymentProgressing)

	if deployment.Status.UpdatedReplicas != *(deployment.Spec.Replicas) ||
		deployment.Status.Replicas != *(deployment.Spec.Replicas) ||
		deployment.Status.AvailableReplicas != *(deployment.Spec.Replicas) ||
		availableCondition.Status != "True" ||
		progressingCondition.Status != "True" {
		return 308, []string{"deployments not ready!"}, nil
	}

	if deployment.Status.ObservedGeneration < deployment.Generation {
		return 308, []string{"observed generation less than generation!"}, nil
	}

	// 发布成功
	return 200, []string{}, nil
}

// GetDeploymentCondition returns the condition with the provided type.
func GetDeploymentCondition(status appsv1.DeploymentStatus, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func GetPodCondition(status corev1.PodStatus, condType corev1.PodConditionType) *corev1.PodCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func GetDeploymentPods(c *gin.Context) {
	appG := app.Gin{C: c}

	var u DeploymentUri

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	deployment, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Get(context.TODO(), u.DeploymentName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	labelSelector := ""
	for key, value := range deployment.Spec.Selector.MatchLabels {
		labelSelector = labelSelector + key + "=" + value + ","
	}
	labelSelector = strings.TrimRight(labelSelector, ",")
	pods, err := k8sClient.ClientV1.CoreV1().Pods(deployment.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", pods)
}

func ForceUpdate(deployment *appsv1.Deployment) {
	if deployment.Spec.Template.Annotations == nil {
		annotations := make(map[string]string)
		annotations["Deployment.UpdateTimestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
		deployment.Spec.Template.Annotations = annotations
	} else {
		deployment.Spec.Template.Annotations["Deployment.UpdateTimestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	}
}
