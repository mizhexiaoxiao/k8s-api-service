package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DestinationRulesQuery struct {
	Namespace string `form:"namespace"`
}

type DestinationRulesUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type DestinationRuleQuery struct {
}

type DestinationRuleUri struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	DrName    string `uri:"drName" binding:"required"`
}

func GetDestinationRules(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q DestinationRulesQuery
		u DestinationRulesUri
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sclient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	istioclient := versionedclient.NewForConfigOrDie(k8sclient.RestConfig)
	drList, err := istioclient.NetworkingV1alpha3().DestinationRules(q.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", drList)
}

func GetDestinationRule(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q DestinationRuleQuery
		u DestinationRuleUri
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sclient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	istioclient := versionedclient.NewForConfigOrDie(k8sclient.RestConfig)
	dr, err := istioclient.NetworkingV1alpha3().DestinationRules(u.Namespace).Get(context.TODO(), u.DrName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", dr)
}

func PostDestinationRule(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u DestinationRulesUri
		q DestinationRulesQuery
		b v1alpha3.DestinationRule
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sclient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	istioclient := versionedclient.NewForConfigOrDie(k8sclient.RestConfig)
	dr, err := istioclient.NetworkingV1alpha3().DestinationRules(q.Namespace).Create(context.TODO(), &b, metav1.CreateOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", dr)
}

func PutDestinationRule(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q DestinationRuleQuery
		u DestinationRuleUri
		b networkingv1alpha3.DestinationRule
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sclient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	istioclient := versionedclient.NewForConfigOrDie(k8sclient.RestConfig)
	dr, err := istioclient.NetworkingV1alpha3().DestinationRules(u.Namespace).Get(context.TODO(), u.DrName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	dr.Spec = b
	_, err = istioclient.NetworkingV1alpha3().DestinationRules(u.Namespace).Update(context.TODO(), dr, metav1.UpdateOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", dr)
}

func DeleteDestinationRule(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u DestinationRuleUri
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sclient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	istioclient := versionedclient.NewForConfigOrDie(k8sclient.RestConfig)

	err = istioclient.NetworkingV1alpha3().DestinationRules(u.Namespace).Delete(context.TODO(), u.DrName, metav1.DeleteOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", nil)
}
