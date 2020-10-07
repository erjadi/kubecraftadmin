package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/sandertv/mcwss/mctype"

	"github.com/sandertv/mcwss/protocol/command"
	"github.com/sandertv/mcwss/protocol/event"

	"github.com/sandertv/mcwss"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var initpos mctype.Position
var initialized bool
var uniqueIDs []string

var agent mcwss.Agent
var namespacesp []mctype.Position

func main() {
	uniqueIDs = make([]string, 0)
	initialized = false
	rand.Seed(86)
	// Create a new server using the default configuration. To use specific configuration, pass a *wss.Config{} in here.
	var c = mcwss.Config{HandlerPattern: "/ws", Address: "0.0.0.0:8000"}
	server := mcwss.NewServer(&c)

	fmt.Println("Listening")

	// Initialize Kube connection
	config, err := clientcmd.BuildConfigFromFlags("", "/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// On first connection
	server.OnConnection(func(player *mcwss.Player) {
		go MOTD(player)
		fmt.Println("Player has entered!")
		player.Exec("time set noon", nil)
		player.Exec("weather clear", nil)

		player.OnTravelled(func(event *event.PlayerTravelled) {
			if !initialized {
				player.Position(func(pos mctype.Position) {
					// Start initialization if you stand on beacon block
					player.Exec("testforblock ~ ~-1 ~ beacon", func(response *command.LocalPlayerName) {
						if response.StatusCode == 0 {

							initpos = pos

							namespacesp = []mctype.Position{
								{pos.X - 11, pos.Y + 5, pos.Z - 11},
								{pos.X - 11, pos.Y + 5, pos.Z - 5},
								{pos.X - 5, pos.Y + 5, pos.Z - 11},
								{pos.X - 5, pos.Y + 5, pos.Z - 5},
							}

							if !initialized {
								initialized = true
								fmt.Println("initialized!")
								Actionbar(player, "Connected to k8s cluster")
								go LoopReconcile(player, clientset)
							}
						}
					})
				})
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
				InitArea(player)
			}

			if (strings.Compare(event.Message, "sync")) == 0 {
				fmt.Println("start syncing")
				go LoopReconcile(player, clientset)
			}

			// Summon 4 animals
			if (strings.Compare(event.Message, "animals")) == 0 {
				fmt.Println("summon animals!")
				Summon(player, initpos, -12+rand.Intn(3), 5, -12+rand.Intn(3), "pig", "test")
				Summon(player, initpos, -12+rand.Intn(3), 5, -6+rand.Intn(3), "chicken", "test")
				Summon(player, initpos, -6+rand.Intn(3), 5, -12+rand.Intn(3), "cow", "test")
				Summon(player, initpos, -6+rand.Intn(3), 5, -6+rand.Intn(3), "horse", "test")
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
