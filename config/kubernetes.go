package config

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Kubernetes struct {
	k8sClient *kubernetes.Clientset
}

func (c *Kubernetes) Client() *kubernetes.Clientset {
	return c.k8sClient
}

func LoadKubernetes() (*Kubernetes, error) {
	inClusterConf, err := rest.InClusterConfig()
	if err != nil {
		log.Print("failed when get in-cluster config")
		return nil, err
	}
	k8sClient, err := kubernetes.NewForConfig(inClusterConf)
	if err != nil {
		return nil, err
	}
	conf := Kubernetes{
		k8sClient: k8sClient,
	}
	return &conf, nil
}
