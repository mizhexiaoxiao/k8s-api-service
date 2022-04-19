package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobsQuery struct {
	Namespace string `form:"namespace"`
	Label     string `form:"label"`
}

type JobsUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type JobUri struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	JobName   string `uri:"jobName" binding:"required"`
}

func GetJobs(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u        JobsUri
		q        JobsQuery
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
		listOpts = metav1.ListOptions{}
	} else {
		listOpts = metav1.ListOptions{LabelSelector: q.Label}
	}
	jobs, err := k8sClient.ClientV1.BatchV1().Jobs(q.Namespace).List(context.TODO(), listOpts)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", jobs)
}

func GetJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u JobUri
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	cronjob, err := k8sClient.ClientV1.BatchV1().Jobs(u.Namespace).Get(context.TODO(), u.JobName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", cronjob)
}

func DeleteJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var u JobUri

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	propagationPolicy := metav1.DeletePropagationBackground
	err = k8sClient.ClientV1.BatchV1().Jobs(u.Namespace).Delete(context.TODO(), u.JobName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", nil)
}
