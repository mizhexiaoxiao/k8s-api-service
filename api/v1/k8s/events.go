package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventsUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type EventsQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Name      string `json:"name" form:"name" binding:"required"`
	Kind      string `json:"kind" form:"kind" binding:"required"`
	Uid       string `json:"uid" form:"uid" binding:"required"`
}

func GetEvents(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u        EventsUri
		q        EventsQuery
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

	clientset, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	listOpts.FieldSelector = fmt.Sprintf(
		"involvedObject.name=%s,involvedObject.namespace=%s,involvedObject.kind=%s,involvedObject.uid=%s",
		q.Name, q.Namespace, q.Kind, q.Uid,
	)
	listOpts.TypeMeta = metav1.TypeMeta{Kind: q.Kind}
	events, err := clientset.CoreV1().Events(q.Namespace).List(context.TODO(), listOpts)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", events)
}
