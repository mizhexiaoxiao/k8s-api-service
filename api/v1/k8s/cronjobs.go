package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJobssQuery struct {
	Namespace string `form:"namespace"`
}

type CronJobsUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type CronJobUri struct {
	Cluster     string `uri:"cluster" binding:"required"`
	Namespace   string `uri:"namespace" binding:"required"`
	CronJobName string `uri:"cronjobName" binding:"required"`
}

type CronJobBody struct {
	Schedule string `json:"schedule" form:"schedule"`
	Suspend  string `json:"suspend" form:"suspend"`
}

func GetCronJobs(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u CronJobsUri
		q CronJobssQuery
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

	cronjobs, err := k8sClient.ClientV1.BatchV1beta1().CronJobs(q.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", cronjobs)
}

func GetCronJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u CronJobUri
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

	cronjob, err := k8sClient.ClientV1.BatchV1beta1().CronJobs(u.Namespace).Get(context.TODO(), u.CronJobName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", cronjob)
}

func PutCronJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		b CronJobBody
		u CronJobUri
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

	cronjob, err := k8sClient.ClientV1.BatchV1beta1().CronJobs(u.Namespace).Get(context.TODO(), u.CronJobName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	if b.Schedule != "" {
		cronjob.Spec.Schedule = b.Schedule
	}

	if b.Suspend != "" {
		suspend, err := strconv.ParseBool(b.Suspend)
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
		cronjob.Spec.Suspend = &suspend
	}

	ucronjob, err := k8sClient.ClientV1.BatchV1beta1().CronJobs(u.Namespace).Update(context.TODO(), cronjob, metav1.UpdateOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "Updated CronJob Successfully", ucronjob)
}

func DeleteCronJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		u CronJobUri
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

	err = k8sClient.ClientV1.BatchV1beta1().CronJobs(u.Namespace).Delete(context.TODO(), u.CronJobName, metav1.DeleteOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "Deleted CronJob Successfully", nil)
}
