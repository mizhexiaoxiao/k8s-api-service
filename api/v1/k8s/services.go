package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServicesQuery struct {
	Namespace string `form:"namespace"`
}

type ServicesUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type ServiceQuery struct {
}

type ServiceUri struct {
	Cluster     string `uri:"cluster" binding:"required"`
	Namespace   string `uri:"namespace" binding:"required"`
	ServiceName string `uri:"serviceName" binding:"required"`
}

func GetServices(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q ServicesQuery
		u ServicesUri
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	clientset, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
	}

	deployments, err := clientset.CoreV1().Services(q.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", deployments)
}

func GetService(c *gin.Context) {
	appG := app.Gin{C: c}
	var u ServiceUri
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	clientset, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
	}

	service, err := clientset.CoreV1().Services(u.Namespace).Get(context.TODO(), u.ServiceName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", service)
}
