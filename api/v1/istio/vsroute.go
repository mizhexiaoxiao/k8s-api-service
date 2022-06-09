package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/istio"
)

type VSHttpRoutesUri struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	VSName    string `uri:"vsName" binding:"required"`
}

type VSHttpRouteUri struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	VSName    string `uri:"vsName" binding:"required"`
	RouteName string `uri:"routeName" binding:"required"`
}

type VSHttpRouteQuery struct {
	APPName string `form:"appName" binding:"required"`
}

func GetVSHttpRoutes(c *gin.Context) {
	var (
		u VSHttpRoutesUri
		q VSHttpRouteQuery
	)
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	istioclient, err := istio.NewIstioClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := istio.NewVSHttpRouteOperation(istioclient, u.Namespace)
	vs, err := operation.List(context.TODO(), u.VSName, q.APPName)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", vs)
}

func GetVSHttpRoute(c *gin.Context) {
	var u VSHttpRouteUri
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	istioclient, err := istio.NewIstioClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := istio.NewVSHttpRouteOperation(istioclient, u.Namespace)
	vs, err := operation.Get(context.TODO(), u.VSName, u.RouteName)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	routes := vs.Spec.Http
	if routes == nil {
		appG.Fail(http.StatusNotFound, fmt.Errorf("VirtualService HTTPRoute %q not found", u.RouteName), nil)
		return
	}
	appG.Success(http.StatusOK, "ok", vs)
}

func AddVSHttpRoute(c *gin.Context) {
	var (
		u VSHttpRoutesUri
		b istio.VSRoute
	)
	appG := app.Gin{C: c}
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	istioclient, err := istio.NewIstioClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := istio.NewVSHttpRouteOperation(istioclient, u.Namespace)
	result, err := operation.Create(context.TODO(), u.VSName, &b)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

func UpdateVSHttpRoute(c *gin.Context) {
	var (
		u VSHttpRouteUri
		b istio.VSRoute
	)
	appG := app.Gin{C: c}
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	istioclient, err := istio.NewIstioClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := istio.NewVSHttpRouteOperation(istioclient, u.Namespace)
	result, err := operation.Update(context.TODO(), u.VSName, u.RouteName, &b)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}

func DeleteVSHttpRoute(c *gin.Context) {
	var (
		u VSHttpRouteUri
		q istio.VSRoute
	)
	appG := app.Gin{C: c}
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindJSON(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	istioclient, err := istio.NewIstioClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	operation := istio.NewVSHttpRouteOperation(istioclient, u.Namespace)
	result, err := operation.Delete(context.TODO(), u.VSName, u.RouteName, &q)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", result)
}
