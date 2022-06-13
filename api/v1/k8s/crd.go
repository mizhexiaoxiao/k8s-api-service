package v1

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"net/http"
)

// GetCRD
// @Summary 获取CRD自定义资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param group path string true "Group"
// @Param version path string true "Version"
// @Param resource path string true "Resource"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/crd/{group}/{version}/{resource}/{namespace}/{name} [get]
func GetCRD(c *gin.Context) {
	appG := app.Gin{C: c}
	param, err := app.GetPathParameterString(c, "cluster", "group", "version", "resource", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	dyn, err := dynamic.NewForConfig(k8sClient.RestConfig)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	crdOperation := k8s.NewCRDOperation(dyn)
	gvk := schema.GroupVersionResource{Group: param["group"], Version: param["version"], Resource: param["resource"]}
	unstructured, err := crdOperation.Get(context.TODO(), gvk, param["namespace"], param["name"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", unstructured.Object)
}

// PostCRD
// @Summary 创建CRD资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param group path string true "Group"
// @Param version path string true "Version"
// @Param resource path string true "Resource"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/crd/{group}/{version}/{resource} [post]
func PostCRD(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		bytes []byte
		data  map[string]interface{}
	)
	param, err := app.GetPathParameterString(c, "cluster", "group", "version", "resource")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	bytes, err = ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	if err = json.Unmarshal(bytes, &data); err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	dyn, err := dynamic.NewForConfig(k8sClient.RestConfig)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	crdOperation := k8s.NewCRDOperation(dyn)
	gvk := schema.GroupVersionResource{Group: param["group"], Version: param["version"], Resource: param["resource"]}
	unstructured, err := crdOperation.Create(context.TODO(), gvk, data)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", unstructured.Object)
}

// PutCRD
// @Summary 更新CRD自定义资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param group path string true "Group"
// @Param version path string true "Version"
// @Param resource path string true "Resource"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/crd/{group}/{version}/{resource}/{namespace}/{name} [put]
func PutCRD(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		bytes []byte
		data  map[string]interface{}
	)
	param, err := app.GetPathParameterString(c, "cluster", "group", "version", "resource", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	bytes, err = ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	if err = json.Unmarshal(bytes, &data); err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	dyn, err := dynamic.NewForConfig(k8sClient.RestConfig)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	crdOperation := k8s.NewCRDOperation(dyn)
	gvk := schema.GroupVersionResource{Group: param["group"], Version: param["version"], Resource: param["resource"]}
	unstructured, err := crdOperation.Update(context.TODO(), gvk, param["namespace"], param["name"], data)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", unstructured.Object)
}

// DeleteCRD
// @Summary 删除CRD自定义资源
// @accept application/json
// @Param cluster path string true "Cluster"
// @Param group path string true "Group"
// @Param version path string true "Version"
// @Param resource path string true "Resource"
// @Param namespace path string true "Namespace"
// @Param name path string true "Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /k8s/{cluster}/crd/{group}/{version}/{resource}/{namespace}/{name} [delete]
func DeleteCRD(c *gin.Context) {
	appG := app.Gin{C: c}
	param, err := app.GetPathParameterString(c, "cluster", "group", "version", "resource", "namespace", "name")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	k8sClient, err := k8s.GetClient(param["cluster"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	dyn, err := dynamic.NewForConfig(k8sClient.RestConfig)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	crdOperation := k8s.NewCRDOperation(dyn)
	gvk := schema.GroupVersionResource{Group: param["group"], Version: param["version"], Resource: param["resource"]}
	err = crdOperation.Delete(context.TODO(), gvk, param["namespace"], param["name"])
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", nil)
}
