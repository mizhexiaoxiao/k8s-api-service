package v1

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mizhexiaoxiao/k8s-api-service/utils"

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

type ExtraPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*ExtraPod `json:"items"`
}

type ExtraPod struct {
	*corev1.Pod
	FormatStatus *k8s.PodFormatStatus `json:"formatStatus"`
}

type ExtraPodResp struct {
	Object *ExtraPod `json:"object"`
	Type   string    `json:"type"`
}

// upgrade websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
	newPodItems := make([]*ExtraPod, len(pods.Items))
	for i := range pods.Items {
		pod := pods.Items[i]
		formatStatus, err := k8s.GetFormatStatus(&pod)
		if err != nil {
			appG.Fail(http.StatusInternalServerError, err, nil)
			return
		}
		newPod := ExtraPod{
			Pod:          &pod,
			FormatStatus: formatStatus,
		}
		newPodItems[i] = &newPod
	}
	newPodList := &ExtraPodList{
		TypeMeta: pods.TypeMeta,
		ListMeta: pods.ListMeta,
		Items:    newPodItems,
	}
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", newPodList)
}

func WatchPods(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		q        PodsQuery
		u        PodsUri
		listOpts metav1.ListOptions
		wg       sync.WaitGroup
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.C.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.C.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer ws.Close()

	if q.Label == "" {
		listOpts = metav1.ListOptions{}
	} else {
		listOpts = metav1.ListOptions{LabelSelector: q.Label}
	}

	k8sClient, err := k8s.GetClient(u.Cluster)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	w, err := k8sClient.ClientV1.CoreV1().Pods(q.Namespace).Watch(context.TODO(), listOpts)
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}

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
		for event := range w.ResultChan() {
			pod, ok := event.Object.(*corev1.Pod)
			formatStatus, err := k8s.GetFormatStatus(pod)
			if err != nil {
				appG.C.AbortWithError(http.StatusInternalServerError, err)
			}
			resp := &ExtraPodResp{
				Object: &ExtraPod{
					Pod:          pod,
					FormatStatus: formatStatus,
				},
				Type: string(event.Type),
			}

			if ok {
				ws.WriteJSON(resp)
			}
		}
	}()

	wg.Wait()
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

func DeletePod(c *gin.Context) {
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

	err = k8sClient.ClientV1.CoreV1().Pods(u.Namespace).Delete(context.TODO(), u.PodName, metav1.DeleteOptions{})

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", nil)

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

func DownloadPodContainerLog(c *gin.Context) {
	appG := app.Gin{C: c}
	params, err := app.GetPathParameterString(appG.C, "cluster", "namespace", "podName", "containerName")
	if err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
	}
	k8sClient, err := k8s.GetClient(params["cluster"])
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	podLogOption := &corev1.PodLogOptions{Container: params["containerName"], LimitBytes: utils.Int64Addr(50 * 1024 * 1024)} // 限制导出最大50Mb
	content, err := k8sClient.ClientV1.CoreV1().Pods(params["namespace"]).GetLogs(params["podName"], podLogOption).Do(context.TODO()).Raw()
	if err != nil {
		appG.C.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	appG.C.Writer.Header().Set("Content-Disposition", fmt.Sprintf("Attachment; Filename=%s-%s.log", params["podName"], params["containerName"]))
	appG.C.Data(http.StatusOK, "Application/OCTET-Stream", content)
}

type WebSSHQuery struct {
	Container string `form:"container" binding:"required"`
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
			Command:   []string{"/bin/sh", "-c", "TERM=xterm-256color; export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c \"/bin/bash\" /dev/null || exec /bin/bash) || exec /bin/sh"},
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
