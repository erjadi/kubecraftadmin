package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/sandertv/mcwss/mctype"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sandertv/mcwss/protocol/command"
	"github.com/sandertv/mcwss/protocol/event"

	"github.com/sandertv/mcwss"
)

var initpos mctype.Position
var initialized bool = false
var playerUniqueIdsMap = make(map[string][]string)
var selectednamespaces []string

var agent mcwss.Agent
var namespacesp []mctype.Position
var playerKubeMap = make(map[string][]string)
var playerEntitiesMap = make(map[string][]string)

// ENV paramaters
var passedNamespaces = os.Getenv("namespaces")
var accessWithinCluster = os.Getenv("accessWithinCluster")

func main() {
	if accessWithinCluster == "" {
		accessWithinCluster = "false"
	}

	initialized = false
	rand.Seed(86)

	clientset, _ := GetClient(accessWithinCluster)

	// Create a new server using the default configuration. To use specific configuration, pass a *wss.Config{} in here.
	var c = mcwss.Config{HandlerPattern: "/ws", Address: "0.0.0.0:8000"}
	server := mcwss.NewServer(&c)

	fmt.Println("Listening on port 8000")

	// On first connection
	server.OnConnection(func(player *mcwss.Player) {
		uniqueIDs := make([]string, 0)
		playerUniqueIdsMap[player.Name()] = uniqueIDs

		//MOTD(player)
		MOTD(player)
		Actionbar(player, "Connected to k8s cluster")

		fmt.Println("Player has entered!")
		player.Exec("time set noon", nil)
		player.Exec("weather clear", nil)
		player.Exec("alwaysday", nil)

		// Provide player with 'equipment'
		player.Exec("give @s diamond_sword", nil)
		player.Exec("give @s tnt 25", nil)
		player.Exec("give @s flint_and_steel", nil)

		playerName := player.Name()
		playerTravelMap := make(map[string]bool)
		playerTravelMap[playerName] = false

		playerInitMap := make(map[string]bool)
		playerInitMap[playerName] = false

		GetPlayerPosition(player)
		SetNamespacesPosition()

		player.OnTravelled(func(event *event.PlayerTravelled) {
			player.Exec("testforblock ~ ~-1 ~ beacon", func(response *command.LocalPlayerName) {
				if response.StatusCode == 0 {
					if !playerInitMap[playerName] {
						playerInitMap[playerName] = true
						fmt.Println("initialized!")

						// Read Namespaces Env - Compile list of selected namespaces
						namespaces, _ := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

						if len(passedNamespaces) > 0 {
							passedNamespacesList := strings.Split(passedNamespaces, ",")
							for _, ns := range namespaces.Items {
								for _, envns := range passedNamespacesList {
									if strings.EqualFold(ns.Name, envns) {
										selectednamespaces = append(selectednamespaces, ns.Name)
									}
								}
							}
							if len(selectednamespaces) < 4 { // if less than 4 specified, select until length is 4
								for _, ns := range namespaces.Items {
									if !Contains(selectednamespaces, ns.Name) {
										selectednamespaces = append(selectednamespaces, ns.Name)
										if len(selectednamespaces) == 4 {
											break
										}
									}
								}
							}
						} else {
							for i := 0; i < 4; i++ {
								selectednamespaces = append(selectednamespaces, namespaces.Items[i].Name)
								fmt.Println("namespace ", selectednamespaces)
							}
						}

						fmt.Println("Selected namespaces: ", selectednamespaces)

						go LoopReconcile(player, clientset)
					}
				}
			})
		})

		// If a mob is killed by the player we do another check which entity is missing
		player.OnMobKilled(func(event *event.MobKilled) {
			fmt.Printf("mobkilled %d\n", event.MobType)
			ReconcileMCtoKubeMob(player, clientset, event.MobType)
		})

		// Set up event handler for commands typed by player
		player.OnPlayerMessage(func(event *event.PlayerMessage) {
			fmt.Println(event.Message)
			if (strings.Compare(event.Message, "detect")) == 0 {
			}

			// Initialize admin area
			if (strings.Compare(event.Message, "init")) == 0 {
				DeleteEntities(player)
				GetPlayerPosition(player)
				SetNamespacesPosition()
				InitArea(player)
			}

			// Force sync if auto-init doesn't work
			if (strings.Compare(event.Message, "sync")) == 0 {
				fmt.Println("start syncing")
				go LoopReconcile(player, clientset)
			}
		})

	})
	server.OnDisconnection(func(player *mcwss.Player) {
		// Called when a player disconnects from the server.
		fmt.Println("Player has disconnected")
	})

	// Run the server. (blocking)
	server.Run()
}
