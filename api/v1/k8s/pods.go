package v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsQuery struct {
	Namespace string `form:"namespace"`
	Label     string `form:"label"`
}

type PodsUri struct {
	Cluster string `uri:"cluster" binding:"required"`
}

type PodQuery struct {
}

type PodUri struct {
	Cluster   string `uri:"cluster" binding:"required"`
	Namespace string `uri:"namespace" binding:"required"`
	PodName   string `uri:"podName" binding:"required"`
}

func GetPods(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q        PodsQuery
		u        PodsUri
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
	if q.Label == "" {
		listOpts = metav1.ListOptions{}
	} else {
		listOpts = metav1.ListOptions{LabelSelector: q.Label}
	}

	clientset, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	pods, err := clientset.CoreV1().Pods(q.Namespace).List(context.TODO(), listOpts)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", pods)
}

func GetPod(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u PodUri
	)
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	clientset, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	pod, err := clientset.CoreV1().Pods(u.Namespace).Get(context.TODO(), u.PodName, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", pod)
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PodLogQuery struct {
	Container  string `form:"container" binding:"required"`
	Follow     string `form:"follow" binding:"required"`
	Previous   string `form:"previous" binding:"required"`
	Timestamps string `form:"timestamps" binding:"required"`
	TailLines  string `form:"tailLines" binding:"required"`
}

func GetPodLog(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u         PodUri
		q         PodLogQuery
		podLogOps corev1.PodLogOptions
	)

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer ws.Close()

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.C.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.C.AbortWithError(http.StatusBadRequest, err)
		return
	}

	follow, _ := strconv.ParseBool(q.Follow)
	previous, _ := strconv.ParseBool(q.Previous)
	timestamps, _ := strconv.ParseBool(q.Timestamps)
	tailLines, _ := strconv.ParseInt(q.TailLines, 10, 64)

	podLogOps = corev1.PodLogOptions{
		Container:  q.Container,
		Follow:     follow,
		Previous:   previous,
		Timestamps: timestamps,
		TailLines:  &tailLines,
	}

	clientset, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	req := clientset.CoreV1().Pods(u.Namespace).GetLogs(u.PodName, &podLogOps)
	readCloser, err := req.Stream(context.TODO())
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer readCloser.Close()
	bufReader := bufio.NewReader(readCloser)

	//cancel context
	ctx, cancel := context.WithCancel(c.Request.Context())
	appG.C.Request = c.Request.WithContext(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	// The goroutine listens to the websocket. When the client cancels,
	// it executes the cancel() callback to close the request context.
	go func() {
		for {
			if _, _, err := ws.NextReader(); err != nil {
				cancel()
				ws.Close()
				wg.Done()
				break
			}
		}
	}()

	go func() {
		for {
			line, _, err := bufReader.ReadLine()
			if err != nil && err != io.EOF {
				ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				break
			}
			if len(line) != 0 {
				err = ws.WriteMessage(websocket.TextMessage, []byte(line))
				if err != nil {
					ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
					break
				}
			}
		}
	}()

	wg.Wait()
}

func HealthCheck(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Success(http.StatusOK, "ok", nil)
	return
}
