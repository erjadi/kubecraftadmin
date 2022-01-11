package main

import (
	"errors"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Contains checks for the occurence of a string in an array of strings
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Remove certain value from an array of strings
func Remove(a []string, x string) []string {
	if Contains(a, x) {
		for i, n := range a {
			if x == n {
				return append(a[:i], a[i+1:]...)
			}
		}
	}
	return a
}

func GetClient(accessWithinCluster string) (*kubernetes.Clientset, error) {
	var client *kubernetes.Clientset
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return client, err
	}

	client, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}

func getKubeConfig(accessWithinCluster string) (error, *rest.Config) {
	var kubeConfig *rest.Config
	var err error
	if accessWithinCluster == "true" {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return err, nil
		}
	} else if accessWithinCluster == "false" {
		var kubeconfigFile string
		if home := homeDir(); home != "" {
			kubeconfigFile = "/.kube" + "/config"
			if _, err := os.Stat(kubeconfigFile); os.IsNotExist(err) {
				return err, nil
			}
		} else {
			return err, nil
		}

		// use the current context in kubeconfig
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
		if err != nil {
			return err, nil
		}
	} else {
		err := errors.New("Parameter AccessWithinCluster not set to true or false")
		return err, nil
	}

	return nil, kubeConfig
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
