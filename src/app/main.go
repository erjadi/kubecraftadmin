package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/sandertv/mcwss/mctype"

	"github.com/sandertv/mcwss/protocol/command"
	"github.com/sandertv/mcwss/protocol/event"

	"github.com/sandertv/mcwss"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var initpos mctype.Position
var initialized bool = false
var uniqueIDs []string
var selectednamespaces []string

var agent mcwss.Agent
var namespacesp []mctype.Position

// ENV paramaters
var passedNamespaces = os.Getenv("namespaces")
var accessWithinCluster = os.Getenv("accessWithinCluster")

func main() {
	if passedNamespaces == "" {
		fmt.Print("The namespaces env parameter was not set (comma separated list of up to 4 namespaces to view in minecraft)!\n")
		os.Exit(1)
	}

	if accessWithinCluster == "" {
		fmt.Print("The accessWithinCluster env parameter was not set (true|false)!\n")
		os.Exit(1)
	}

	uniqueIDs = make([]string, 0)
	initialized = false
	rand.Seed(86)

	clientset, _ := GetClient(accessWithinCluster)

	// Create a new server using the default configuration. To use specific configuration, pass a *wss.Config{} in here.
	var c = mcwss.Config{HandlerPattern: "/ws", Address: "0.0.0.0:8000"}
	server := mcwss.NewServer(&c)

	fmt.Println("Listening on port 8000")

	// On first connection
	server.OnConnection(func(player *mcwss.Player) {
		go MOTD(player)
		fmt.Println("Player has entered!")
		player.Exec("time set noon", nil)
		player.Exec("weather clear", nil)
		player.Exec("alwaysday", nil)

		// Provide player with 'equipment'
		player.Exec("give @s diamond_sword", nil)
		player.Exec("give @s tnt 25", nil)
		player.Exec("give @s flint_and_steel", nil)

		fmt.Println("Selected namespaces: ", selectednamespaces)

		player.OnTravelled(func(event *event.PlayerTravelled) {
			if !initialized {
				//initpos = GetPlayerPosition(player)
				//namespacesp = GetNamespacesPosition(initpos)
				var x float64
				var y float64
				var z float64

				player.Exec("tp ~~~", func(response map[string]interface{}) {
					if destination, ok := response["destination"]; ok {
						xString := fmt.Sprintf("%v", destination.(interface{}).(map[string]interface{})["x"])
						x, _ = strconv.ParseFloat(xString, 64)

						yString := fmt.Sprintf("%v", destination.(interface{}).(map[string]interface{})["y"])
						y, _ = strconv.ParseFloat(yString, 64)

						zString := fmt.Sprintf("%v", destination.(interface{}).(map[string]interface{})["z"])
						z, _ = strconv.ParseFloat(zString, 64)
					}

					initpos.X = x
					initpos.Y = y
					initpos.Z = z

				})

				namespacesp = []mctype.Position{
					{X: initpos.X - 11, Y: initpos.Y + 5, Z: initpos.Z - 11},
					{X: initpos.X - 11, Y: initpos.Y + 5, Z: initpos.Z - 5},
					{X: initpos.X - 5, Y: initpos.Y + 5, Z: initpos.Z - 11},
					{X: initpos.X - 5, Y: initpos.Y + 5, Z: initpos.Z - 5},
				}

				player.Exec("testforblock ~ ~-1 ~ beacon", func(response *command.LocalPlayerName) {
					if response.StatusCode == 0 {
						if !initialized {
							initialized = true
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

							Actionbar(player, "Connected to k8s cluster")
							go LoopReconcile(player, clientset)
						}
					}
				})
				//})
			}
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
				//player.Position(func(pos mctype.Position) {
				// Start initialization if you stand on beacon block
				//initpos = GetPlayerPosition(player)
				//namespacesp = GetNamespacesPosition(initpos)

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
