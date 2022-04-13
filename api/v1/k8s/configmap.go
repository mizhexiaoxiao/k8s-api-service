package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	"github.com/mizhexiaoxiao/k8s-api-service/models/metadata"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

// PostConfigmap
// @Summary 创建Configmap资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/configmaps [post]
func PostConfigmap(c *gin.Context) {
	appG := app.Gin{C: c}
	var configMap v1.ConfigMap

	param, err := app.GetPathParameterString(c, "cluster")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBind(&configMap); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	configMapOperation := k8s.NewConfigmapOperation(k8sClient.ClientV1)
	result, err := configMapOperation.Create(context.TODO(), &configMap)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// GetConfigmapList
// @Summary 获取Configmap资源列表
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param param query metadata.CommonQueryParameter true "LabelSelector"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/configmaps [get]
func GetConfigmapList(c *gin.Context) {
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

	configMapOperation := k8s.NewConfigmapOperation(k8sClient.ClientV1)
	result, err := configMapOperation.List(context.TODO(), queryParam)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// GetConfigmap
// @Summary 获取Configmap资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/configmaps/{namespace}/{name} [get]
func GetConfigmap(c *gin.Context) {
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
	configMapOperation := k8s.NewConfigmapOperation(k8sClient.ClientV1)
	configMap, err := configMapOperation.Get(context.TODO(), param["namespace"], param["name"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", configMap)
}

// PutConfigmap
// @Summary 更新Configmap资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/configmaps/{namespace}/{name} [put]
func PutConfigmap(c *gin.Context) {
	appG := app.Gin{C: c}
	var configMap v1.ConfigMap
	param, err := app.GetPathParameterString(c, "cluster", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBind(&configMap); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	configMapOperation := k8s.NewConfigmapOperation(k8sClient.ClientV1)
	result, err := configMapOperation.Update(context.TODO(), param["namespace"], param["name"], &configMap)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

// DeleteConfigmap
// @Summary 删除Configmap资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/configmaps/{namespace}/{name} [delete]
func DeleteConfigmap(c *gin.Context) {
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

	configMapOperation := k8s.NewConfigmapOperation(k8sClient.ClientV1)
	err = configMapOperation.Delete(context.TODO(), param["namespace"], param["name"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", nil)
}
