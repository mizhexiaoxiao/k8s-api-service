package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/cluster"
	"github.com/mizhexiaoxiao/k8s-api-service/models"
)

func PostCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		b models.ClusterModel
	)
	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := cluster.Create(b); err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "Created Successfully", nil)
}

func PutCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u app.GetById
		b models.ClusterModel
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindJSON(&b); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := cluster.Update(u.ID, b); err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "Updated Successfully", nil)
}

func ListCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		pageInfo app.PageInfo
	)
	if err := appG.C.ShouldBindQuery(&pageInfo); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	res, err := cluster.List(pageInfo)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	count, err := cluster.Count()
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.SuccessExtra(count, pageInfo.Page, pageInfo.PageSize, http.StatusOK, "ok", res)
}

func GetCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u app.GetById
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	res, err := cluster.Get(u.ID)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", res)
}

func DeleteCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		idInfo app.GetById
	)
	if err := appG.C.ShouldBindUri(&idInfo); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	if err := cluster.Delete(idInfo.ID); err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "Deleted Successfully", nil)
}

func TestConnectCluster(c *gin.Context) {

}
