package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentsQuery struct {
	Namespace string `form:"namespace"`
	Label     string `form:"label"`
}

type DeploymentActionQuery struct {
	Action string `form:"action" binding:"required"`
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

var APIVersion = "apps/v1"
var Kind = "Deployment"

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
	deployment.TypeMeta.APIVersion = APIVersion
	deployment.TypeMeta.Kind = Kind
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
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/deployments [post]
func PostDeployment(c *gin.Context) {
	appG := app.Gin{C: c}
	var deployment appsv1.Deployment

	cluster := appG.C.Param("cluster")
	if cluster == "" {
		appG.Fail(http.StatusBadRequest, errors.New("cluster param not valid"), nil)
		return
	}
	if err := appG.C.ShouldBind(&deployment); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	operation := k8s.NewDeploymentOperation(k8sClient.ClientV1)
	result, err := operation.Create(context.TODO(), &deployment)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

func DeploymentDoAction(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u DeploymentUri
		q DeploymentActionQuery
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

	deployment, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Get(context.TODO(), u.DeploymentName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	switch q.Action {
	case "redeploy":
		if deployment.Spec.Paused {
			appG.Fail(http.StatusInternalServerError, errors.New("can't restart paused deployment (run rollout resume first)"), nil)
			return
		}
		if deployment.Spec.Template.ObjectMeta.Annotations == nil {
			deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
		}
		deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().String()
		_, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
	case "pause":
		if deployment.Spec.Paused {
			appG.Fail(http.StatusInternalServerError, errors.New("deployment is already paused"), nil)
			return
		}
		deployment.Spec.Paused = true
		_, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
	case "resume":
		if !deployment.Spec.Paused {
			appG.Fail(http.StatusInternalServerError, errors.New("deployment is not paused"), nil)
			return
		}
		deployment.Spec.Paused = false
		_, err := k8sClient.ClientV1.AppsV1().Deployments(u.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
	default:
		appG.Fail(http.StatusBadRequest, errors.New("Invalid parameter, must be redeploy|pause|resume"), nil)
		return
	}

	appG.Success(http.StatusOK, "ok", nil)

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

// PatchDeployment
// @Summary 批量更新deployment
// @Produce  json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "DeploymentName"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/deployments/{namespace}/{deploymentName} [patch]
func PatchDeployment(c *gin.Context) {
	appG := app.Gin{C: c}
	var deployment appsv1.Deployment

	params, err := app.GetPathParameterString(c, "cluster", "namespace", "deploymentName")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBind(&deployment); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(params["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	operation := k8s.NewDeploymentOperation(k8sClient.ClientV1)
	result, err := operation.Update(context.TODO(), params["namespace"], params["deploymentName"], &deployment)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
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
		status, reasons, err := getDeploymentStatus(deployment)
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, reasons)
			return
		}
		if status == http.StatusOK {
			appG.Success(http.StatusOK, reasons, nil)
			return
		}
		if status == http.StatusPermanentRedirect {
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
			status, reasons, err := getDeploymentStatus(&deployment)
			if err != nil {
				appG.Fail(http.StatusInternalServerError, err, nil)
				return
			}
			if status == http.StatusOK {
				appG.Success(http.StatusOK, reasons, nil)
				return
			}
			if status == http.StatusPermanentRedirect {
				appG.Fail(http.StatusPermanentRedirect, errors.New("retry"), reasons)
				return
			}
		}
		appG.Success(http.StatusOK, "ok", nil)
		return
	}
}

const (
	TimedOutReason = "ProgressDeadlineExceeded"
)

func getDeploymentStatus(deployment *appsv1.Deployment) (status int, reasons string, err error) {
	if deployment.Generation <= deployment.Status.ObservedGeneration {
		cond := GetDeploymentCondition(deployment.Status, appsv1.DeploymentProgressing)
		if cond != nil && cond.Reason == TimedOutReason {
			return http.StatusInternalServerError, "", fmt.Errorf("deployment %s exceeded its progress deadline", deployment.Name)
		}
		if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
			return http.StatusPermanentRedirect, fmt.Sprintf("Waiting for deployment %s rollout to finish: %d out of %d new replicas have been updated...", deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas), nil
		}
		if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
			return http.StatusPermanentRedirect, fmt.Sprintf("Waiting for deployment %s rollout to finish: %d old replicas are pending termination...", deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas), nil
		}
		if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
			return http.StatusPermanentRedirect, fmt.Sprintf("Waiting for deployment %s rollout to finish: %d of %d updated replicas are available...", deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas), nil
		}
		return http.StatusOK, fmt.Sprintf("deployment %s successfully rolled out", deployment.Name), nil
	}
	return http.StatusPermanentRedirect, fmt.Sprintf("Waiting for deployment spec update to be observed..."), nil
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
