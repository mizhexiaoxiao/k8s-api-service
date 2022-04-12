package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	"github.com/mizhexiaoxiao/k8s-api-service/models/metadata"
	v1 "k8s.io/api/autoscaling/v1"
)

// PostHorizontalPodAutoScalers
// @Summary 创建弹性伸缩资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/horizontalpodautoscalers [post]
func PostHorizontalPodAutoScalers(c *gin.Context) {
	appG := app.Gin{C: c}
	var scaler v1.HorizontalPodAutoscaler

	param, err := app.GetPathParameterString(c, "cluster")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBind(&scaler); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	operation := k8s.NewHorizontalPodAutoScalerOperation(k8sClient.ClientV1)
	result, err := operation.Create(context.TODO(), &scaler)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// GetHorizontalPodAutoScalerList
// @Summary 获取弹性伸缩资源列表
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace query string true "Namespace"
// @Param param query metadata.CommonQueryParameter true "LabelSelector"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/horizontalpodautoscalers [get]
func GetHorizontalPodAutoScalerList(c *gin.Context) {
	appG := app.Gin{C: c}
	var queryParam metadata.CommonQueryParameter
	pathParam, err := app.GetPathParameterString(c, "cluster")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&queryParam); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(pathParam["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	operation := k8s.NewHorizontalPodAutoScalerOperation(k8sClient.ClientV1)
	result, err := operation.List(context.TODO(), queryParam)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// GetHorizontalPodAutoScaler
// @Summary 获取弹性伸缩资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/horizontalpodautoscalers/{namespace}/{name} [get]
func GetHorizontalPodAutoScaler(c *gin.Context) {
	appG := app.Gin{C: c}
	param, err := app.GetPathParameterString(c, "cluster", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := k8s.NewHorizontalPodAutoScalerOperation(k8sClient.ClientV1)
	result, err := operation.Get(context.TODO(), param["namespace"], param["name"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// PutHorizontalPodAutoScaler
// @Summary 更新弹性伸缩资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/horizontalpodautoscalers/{namespace}/{name} [put]
func PutHorizontalPodAutoScaler(c *gin.Context) {
	appG := app.Gin{C: c}
	var scaler v1.HorizontalPodAutoscaler
	param, err := app.GetPathParameterString(c, "cluster", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBind(&scaler); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := k8s.NewHorizontalPodAutoScalerOperation(k8sClient.ClientV1)
	result, err := operation.Update(context.TODO(), param["namespace"], param["name"], &scaler)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// DeleteHorizontalPodAutoScaler
// @Summary 删除弹性伸缩资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/horizontalpodautoscalers/{namespace}/{name} [delete]
func DeleteHorizontalPodAutoScaler(c *gin.Context) {
	appG := app.Gin{C: c}
	param, err := app.GetPathParameterString(c, "cluster", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	operation := k8s.NewHorizontalPodAutoScalerOperation(k8sClient.ClientV1)
	err = operation.Delete(context.TODO(), param["namespace"], param["name"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", nil)
}
