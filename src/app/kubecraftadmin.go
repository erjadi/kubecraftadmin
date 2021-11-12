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
var wolf bool = false

// ReconcileKubetoMC queries Kubernetes cluster for resources and removes / spawns entities accordingly in Minecraft
func ReconcileKubetoMC(p *mcwss.Player, clientset *kubernetes.Clientset) {
	synclock.Lock()
	p.Exec("testfor @e", func(response map[string]interface{}) {

		kubeentities := make([]string, 0)

		victims := fmt.Sprintf("%s", response["victim"])
		mcentities := strings.Fields(victims[1 : len(victims)-1])

		// Get all Kube entities per namespace
		for i, ns := range selectednamespaces {

			// Summon fox to eat the animals
			//selectedNs, _ := clientset.CoreV1().Namespaces().Get(ns)

			pods, _ := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
			services, _ := clientset.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
			deployments, _ := clientset.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
			rc, _ := clientset.AppsV1().ReplicaSets(ns).List(context.TODO(), metav1.ListOptions{})
			statefulset, _ := clientset.AppsV1().StatefulSets(ns).List(context.TODO(), metav1.ListOptions{})
			daemonset, _ := clientset.AppsV1().DaemonSets(ns).List(context.TODO(), metav1.ListOptions{})

			for _, pod := range pods.Items {
				//if !wolf {
				//	Summonpos(p, clientset, namespacesp[i], "wolf", "BigBadWolf")
				//	wolf = true
				//}
				kubeentities = append(kubeentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name)) {
					if pod.Status.Phase == v1.PodRunning {
						Summonpos(p, clientset, namespacesp[i], "sheep", fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name))
					}
				}
			}
			for _, deployment := range deployments.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name)) {
					fmt.Printf("Summoning %s:deployment:%s\n", deployment.Namespace, deployment.Name)
					Summonpos(p, clientset, namespacesp[i], "horse", fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name))
				}
			}
			for _, rcontr := range rc.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name)) {
					Summonpos(p, clientset, namespacesp[i], "cow", fmt.Sprintf("%s:replicaset:%s", rcontr.Namespace, rcontr.Name))
				}
			}
			for _, service := range services.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name)) {
					Summonpos(p, clientset, namespacesp[i], "chicken", fmt.Sprintf("%s:service:%s", service.Namespace, service.Name))
				}
			}
			for _, ss := range statefulset.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:statefulset:%s", ss.Namespace, ss.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:statefulset:%s", ss.Namespace, ss.Name)) {
					Summonpos(p, clientset, namespacesp[i], "pig", fmt.Sprintf("%s:statefulset:%s", ss.Namespace, ss.Name))
				}
			}
			for _, ds := range daemonset.Items {
				kubeentities = append(kubeentities, fmt.Sprintf("%s:daemonset:%s", ds.Namespace, ds.Name))
				if !Contains(mcentities, fmt.Sprintf("%s:daemonset:%s", ds.Namespace, ds.Name)) {
					Summonpos(p, clientset, namespacesp[i], "goat", fmt.Sprintf("%s:daemonset:%s", ds.Namespace, ds.Name))
				}
			}
		}

		// Delete entities

		for _, entity := range mcentities {
			if !Contains(kubeentities, entity) {
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=horse]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=sheep]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=goat]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=chicken]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=cow]", entity), nil)
				p.Exec(fmt.Sprintf("kill @e[name=%s,type=pig]", entity), nil)
				uniqueIDs = Remove(uniqueIDs, entity)
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
	fmt.Println("mob type ", mobType)
	if mobType == 13 { // delete pod
		p.Exec("testfor @e[type=sheep]", func(response map[string]interface{}) {

			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			for _, ns := range selectednamespaces {
				// Get all Kube entities for selected namespace
				pods, _ := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})

				for _, pod := range pods.Items {
					if !Contains(mcentities, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name)) {
						fmt.Printf(fmt.Sprintf("Player killed %s:pod:%s\n", pod.Namespace, pod.Name))
						clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
						uniqueIDs = Remove(uniqueIDs, fmt.Sprintf("%s:pod:%s", pod.Namespace, pod.Name))
					}
				}
			}
		})
	}
	if mobType == 12 { // delete pod
		p.Exec("testfor @e[type=pig]", func(response map[string]interface{}) {

			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			for _, ns := range selectednamespaces {
				// Get all Kube entities for selected namespace
				statefulset, _ := clientset.AppsV1().StatefulSets(ns).List(context.TODO(), metav1.ListOptions{})

				for _, ss := range statefulset.Items {
					if !Contains(mcentities, fmt.Sprintf("%s:statefulset:%s", ss.Namespace, ss.Name)) {
						fmt.Printf(fmt.Sprintf("Player killed %s:statefulset:%s\n", ss.Namespace, ss.Name))
						//clientset.CoreV1().StatefulSets(ss.Namespace).Delete(context.TODO(), ss.Name, metav1.DeleteOptions{})
						uniqueIDs = Remove(uniqueIDs, fmt.Sprintf("%s:statefulset:%s", ss.Namespace, ss.Name))
					}
				}
			}
		})
	}
	if mobType == 10 { // delete service
		p.Exec("testfor @e[type=chicken]", func(response map[string]interface{}) {
			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			for _, ns := range selectednamespaces {
				// Get all Kube entities for selected namespace
				services, _ := clientset.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})

				for _, service := range services.Items {
					if !Contains(mcentities, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name)) {
						fmt.Printf(fmt.Sprintf("Player killed %s:service:%s!\n", service.Namespace, service.Name))
						//clientset.CoreV1().Services(service.Namespace).Delete(context.TODO(), service.Name, metav1.DeleteOptions{})
						uniqueIDs = Remove(uniqueIDs, fmt.Sprintf("%s:service:%s", service.Namespace, service.Name))
					}
				}
			}

		})
	}
	if mobType == 11 { // delete replicaset
		p.Exec("testfor @e[type=cow]", func(response map[string]interface{}) {

			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			for _, ns := range selectednamespaces {
				// Get all Kube entities for selected namespace
				rcs, _ := clientset.AppsV1().ReplicaSets(ns).List(context.TODO(), metav1.ListOptions{})

				for _, rc := range rcs.Items {
					if !Contains(mcentities, fmt.Sprintf("%s:replicaset:%s", rc.Namespace, rc.Name)) {
						fmt.Printf(fmt.Sprintf("Player killed %s:replicaset:%s\n", rc.Namespace, rc.Name))
						//clientset.AppsV1().ReplicaSets(rc.Namespace).Delete(context.TODO(), rc.Name, metav1.DeleteOptions{})
						uniqueIDs = Remove(uniqueIDs, fmt.Sprintf("%s:replicaset:%s", rc.Namespace, rc.Name))
					}
				}
			}
		})
	}

	if mobType == 128 { // delete deployment
		p.Exec("testfor @e[type=goat]", func(response map[string]interface{}) {

			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			for _, ns := range selectednamespaces {
				// Get all Kube entities for selected namespace
				daemonset, _ := clientset.AppsV1().DaemonSets(ns).List(context.TODO(), metav1.ListOptions{})

				for _, ds := range daemonset.Items {
					if !Contains(mcentities, fmt.Sprintf("%s:daemonset:%s", ds.Namespace, ds.Name)) {
						fmt.Printf(fmt.Sprintf("Player killed %s:daemonset:%s\n", ds.Namespace, ds.Name))
						//clientset.AppsV1().DaemonSets(ds.Namespace).Delete(context.TODO(), ds.Name, metav1.DeleteOptions{})
						// Remove deployment from uniqueIDs to allow recreation
						uniqueIDs = Remove(uniqueIDs, fmt.Sprintf("%s:daemonset:%s", ds.Namespace, ds.Name))
					}
				}
			}
		})
	}

	if mobType == 23 { // delete deployment
		p.Exec("testfor @e[type=horse]", func(response map[string]interface{}) {

			victims := fmt.Sprintf("%s", response["victim"])
			mcentities := strings.Fields(victims[1 : len(victims)-1])

			for _, ns := range selectednamespaces {
				// Get all Kube entities for selected namespace
				deployments, _ := clientset.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})

				for _, deployment := range deployments.Items {
					if !Contains(mcentities, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name)) {
						fmt.Printf(fmt.Sprintf("Player killed %s:deployment:%s\n", deployment.Namespace, deployment.Name))
						//clientset.AppsV1().Deployments(deployment.Namespace).Delete(context.TODO(), deployment.Name, metav1.DeleteOptions{})
						// Remove deployment from uniqueIDs to allow recreation
						uniqueIDs = Remove(uniqueIDs, fmt.Sprintf("%s:deployment:%s", deployment.Namespace, deployment.Name))
					}
				}
			}
		})
	}
}
