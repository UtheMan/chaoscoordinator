package service

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//func SetConfigs() *kubernetes.Clientset {
//	var kubeconfig *string
//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
//	flag.Parse()
//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil {
//		panic(err.Error())
//	}
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		panic(err.Error())
//	}
//	return clientset
//}
func SetConfigs() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset
}
