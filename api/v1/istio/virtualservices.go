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

type VirtualServicesQuery struct {
	Namespace string `form:"namespace"`
}

type VirtualServicesUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type VirtualServiceQuery struct {
}

type VirtualServiceUri struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	VSName    string `uri:"vsName" binding:"required"`
}

func GetVirtualServices(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q VirtualServicesQuery
		u VirtualServicesUri
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
	vsList, err := istioclient.NetworkingV1alpha3().VirtualServices(q.Namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", vsList)
}

func GetVirtualService(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q VirtualServiceQuery
		u VirtualServiceUri
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
	vs, err := istioclient.NetworkingV1alpha3().VirtualServices(u.Namespace).Get(context.TODO(), u.VSName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", vs)
}

func PostVirtualService(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u VirtualServicesUri
		q VirtualServicesQuery
		b v1alpha3.VirtualService
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
	vs, err := istioclient.NetworkingV1alpha3().VirtualServices(q.Namespace).Create(context.TODO(), &b, metav1.CreateOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", vs)
}

func PutVirtualService(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q VirtualServiceQuery
		u VirtualServiceUri
		b networkingv1alpha3.VirtualService
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
	vs, err := istioclient.NetworkingV1alpha3().VirtualServices(u.Namespace).Get(context.TODO(), u.VSName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	vs.Spec = b
	_, err = istioclient.NetworkingV1alpha3().VirtualServices(u.Namespace).Update(context.TODO(), vs, metav1.UpdateOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", nil)
}

func DeleteVirtualService(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u VirtualServiceUri
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

	err = istioclient.NetworkingV1alpha3().VirtualServices(u.Namespace).Delete(context.TODO(), u.VSName, metav1.DeleteOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", nil)
}
