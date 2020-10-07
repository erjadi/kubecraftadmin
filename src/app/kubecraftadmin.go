package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sandertv/mcwss"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var synclock sync.Mutex

// ReconcileKubetoMC queries Kubernetes cluster for resources and removes / spawns entities accordingly in Minecraft
func ReconcileKubetoMC(p *mcwss.Player, clientset *kubernetes.Clientset) {

	synclock.Lock()
	p.Exec("testfor @e", func(response map[string]interface{}) {

		kubeentities := make([]string, 0)

		victims := fmt.Sprintf("%s", response["victim"])
		mcentities := strings.Fields(victims[1 : len(victims)-1])

		// Get all Kube entities per namespace
		namespaces, _ := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		for i, ns := range namespaces.Items {

			pods, _ := clientset.CoreV1().Pods(ns.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
			services, _ := clientset.CoreV1().Services(ns.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
			deployments, _ := clientset.AppsV1().Deployments(ns.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})
			rc, _ := clientset.AppsV1().ReplicaSets(ns.ObjectMeta.Name).List(context.TODO(), metav1.ListOptions{})

			for _, pod := range pods.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name)) {
					if pod.Status.Phase == v1.PodRunning {
						Summonpos(p, namespacesp[i], "pig", fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name))
					}
				}
			}
			for _, deployment := range deployments.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name)) {
					fmt.Printf("Summoning %s:deployment%s\n", deployment.Namespace, deployment.Name)
					Summonpos(p, namespacesp[i], "horse", fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name))
				}
			}
			for _, rcontr := range rc.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name)) {
					Summonpos(p, namespacesp[i], "cow", fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name))
				}
			}
			for _, service := range services.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name)) {
					Summonpos(p, namespacesp[i], "chicken", fmt.Sprintf("%s:service:%s", service.Namespace, service.Name))
				}
			}
		}

		// Delete entities

		for _, entity := range mcentities {
			if !Contains(kubeentities, entity) {
				// fmt.Printf("Kube side kill %s\n", entity)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=horse]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=pig]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=chicken]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=cow]", entity), nil)
				time.Sleep(50 * time.Millisecond)
			}
		}

		synclock.Unlock()
	})
}

// LoopReconcile will run ReconcileKubetoMC every second
func LoopReconcile(p *mcwss.Player, clientset *kubernetes.Clientset) {
	for {
		ReconcileKubetoMC(p, clientset)
		time.Sleep(1 * time.Second)
	}
}

// ReconcileMCtoKubeMob will delete a specific resource from Kubernetes based on the entities found in Minecraft. Typically run after mob event.
func ReconcileMCtoKubeMob(p *mcwss.Player, clientset *kubernetes.Clientset, mobType int) {
	if mobType == 12 {
		p.Exec("testfor @e[type=pig]", func(response map[string]interface{}) {

			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			// Get all Kube entities
			pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			for _, pod := range pods.Items {
				if !Contains(mcentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name)) {
					fmt.Printf(fmt.Sprintf("Kill %s:pod:%s!!\n", pod.Namespace, pod.Name))
					clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				}
			}
		})
	}
}

// ReconcileMCtoKube DEPRECATED: will query all entities in Minecraft and removes kubernetes resources
func ReconcileMCtoKube(p *mcwss.Player, clientset *kubernetes.Clientset) {

	fmt.Println("MCtoKube")
	p.Exec("testfor @e", func(response map[string]interface{}) {

		victims := fmt.Sprintf("%s", response["victim"])
		mcentities := strings.Fields(victims[1 : len(victims)-1])

		// Get all Kube entities
		pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		services, _ := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		deployments, _ := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		rc, _ := clientset.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})

		for _, pod := range pods.Items {
			if strings.Compare(pod.Namespace, "default") == 0 {
				if !Contains(mcentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name)) {
					fmt.Printf(fmt.Sprintf("Kill %s:pod:%s!!\n", pod.Namespace, pod.Name))
					// freshlydeleted.Put(fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name), "locked")
					// fmt.Printf("%s locked\n", (fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name)))
					clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				}
			}
		}
		for _, deployment := range deployments.Items {
			if !Contains(mcentities, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name)) {
				fmt.Printf(fmt.Sprintf("kill %s:deployment:%s\n", deployment.Namespace, deployment.Name))
			}
		}
		for _, rcontr := range rc.Items {
			if !Contains(mcentities, fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name)) {
				fmt.Printf(fmt.Sprintf("kill %s:replicaset:%s\n", rcontr.Namespace, rcontr.Name))
			}
		}
		for _, service := range services.Items {
			if !Contains(mcentities, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name)) {
				fmt.Printf(fmt.Sprintf("kill %s:service:%s\n", service.Namespace, service.Name))
			}
		}
	})
}
