package k8s

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetK8sClient() (*kubernetes.Clientset, error) {
	kubeConfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	// Use in-cluster config if kubeconfig path is not specified'
	var (
		config *rest.Config
		err    error
	)
	if *kubeConfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
