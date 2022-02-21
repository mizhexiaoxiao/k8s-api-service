package v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
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

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	pods, err := k8sClient.ClientV1.CoreV1().Pods(q.Namespace).List(context.TODO(), listOpts)
	for i := 0; i < len(pods.Items); i++ {
		pods.Items[i].CreationTimestamp = metav1.NewTime(pods.Items[i].CreationTimestamp.Add(8 * time.Hour))
	}

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

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	pod, err := k8sClient.ClientV1.CoreV1().Pods(u.Namespace).Get(context.TODO(), u.PodName, metav1.GetOptions{})
	pod.CreationTimestamp = metav1.NewTime(pod.CreationTimestamp.Add(8 * time.Hour))

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", pod)
}

type PodLogQuery struct {
	Container  string `form:"container" binding:"required"`
	Follow     string `form:"follow" binding:"required"`
	Previous   string `form:"previous" binding:"required"`
	Timestamps string `form:"timestamps" binding:"required"`
	TailLines  string `form:"tailLines" binding:"required"`
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetPodLog(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u         PodUri
		q         PodLogQuery
		podLogOps corev1.PodLogOptions
		wg        sync.WaitGroup
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

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	req := k8sClient.ClientV1.CoreV1().Pods(u.Namespace).GetLogs(u.PodName, &podLogOps)
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

type WebSSHQuery struct {
	Container string `form:"container" binding:"required"`
	Command   string `form:"command" binding:"required"`
}

func PodWebSSH(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		u PodUri
		q WebSSHQuery
		t *WebTerminal
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

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sshReq := k8sClient.ClientV1.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(u.PodName).
		Namespace(u.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: q.Container,
			Command:   []string{q.Command},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	t = &WebTerminal{
		wsConn:   ws,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
		tty:      true,
	}

	// create connection to container
	// sshReq.URL() => https://xxxxxx/api/v1/namespaces/dev/pods/front-tracing-go-58d8dfb599-l4x94/exec?command=bash&container=app&stderr=true&stdin=true&stdout=true&tty=true
	executor, err := remotecommand.NewSPDYExecutor(k8sClient.RestConfig, "POST", sshReq.URL())
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
	}

	// Data flow processing callback between configuration and container
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             t.Stdin(),
		Stdout:            t.Stdout(),
		Stderr:            t.Stderr(),
		Tty:               t.Tty(),
		TerminalSizeQueue: t,
	})

	if err != nil {
		t.wsConn.Close()
		appG.C.AbortWithError(http.StatusInternalServerError, err)
	}

	return
}
