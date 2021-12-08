package k8s

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mizhexiaoxiao/k8s-api-service/models"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

var k8sClients = &sync.Map{} //请求map

func GetClient(clusterName string) (*kubernetes.Clientset, error) {
	var (
		cluster models.Cluster
		context clientcmdapiv1.Config
		err     error
	)
	client, ok := k8sClients.Load(clusterName)
	if ok {
		return client.(*kubernetes.Clientset), nil
	}

	err = models.DB.Model(&models.ClusterModel{}).Where("name = ?", clusterName).First(&cluster).Error
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cluster.Context, &context)
	if err != nil {
		return nil, err
	}
	clientset, err := BuildClient(cluster.Name, context)
	if err != nil {
		return nil, err
	}
	k8sClients.Store(clusterName, clientset)
	return clientset, nil
}

// func GetLocalClient(clusterID string) (*kubernetes.Clientset, error) {
// 	client, ok := k8sClients.Load(clusterID)
// 	if ok {
// 		// logging.Info("从缓存中得到", clusterID, "的连接")
// 		return client.(*kubernetes.Clientset), nil
// 	}
// 	// logging.Info("为", clusterID, "创建连接")

// 	// 从本地读取配置
// 	var kubeconfig *string
// 	if home := homedir.HomeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()
// 	// use the current context in kubeconfig
// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//clientcmd.NewClientConfigFromBytes()

// 	// create the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// logging.Info(clusterID, "连接创建成功")
// 	fmt.Println(clusterID, "连接成功")
// 	k8sClients.Store(clusterID, clientset)
// 	return clientset, nil
// }

const (
	// High enough QPS to fit all expected use cases.
	defaultQPS = 1e6
	// High enough Burst to fit all expected use cases.
	defaultBurst = 1e6
	// full resyc cache resource time
	defaultResyncPeriod = 30 * time.Second
)

func BuildClient(server string, configV1 clientcmdapiv1.Config) (*kubernetes.Clientset, error) {
	configObject, err := clientcmdlatest.Scheme.ConvertToVersion(&configV1, clientcmdapi.SchemeGroupVersion)
	configInternal := configObject.(*clientcmdapi.Config)

	clientConfig, err := clientcmd.NewDefaultClientConfig(*configInternal,
		&clientcmd.ConfigOverrides{
			ClusterDefaults: clientcmdapi.Cluster{Server: server},
		}).ClientConfig()

	if err != nil {
		return nil, err
	}

	clientConfig.QPS = defaultQPS
	clientConfig.Burst = defaultBurst

	clientSet, err := kubernetes.NewForConfig(clientConfig)

	if err != nil {
		return nil, err
	}

	return clientSet, nil
}
