package main

import (
	"fmt"
	"math/rand"
	"strconv"

	//	"math/rand"
	"time"

	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/mctype"
	"github.com/sandertv/mcwss/protocol/command"
	"k8s.io/client-go/kubernetes"
)

// PlayerFill will fill the playing area with blocktype, coordinates are relative to the player position
func PlayerFill(p *mcwss.Player, x1 int, y1 int, z1 int, x2 int, y2 int, z2 int, blocktype string) {
	p.Exec(fmt.Sprintf("fill ~%d ~%d ~%d ~%d ~%d ~%d %s", x1, y1, z1, x2, y2, z2, blocktype), nil)
}

// Fill will fill the playing area with blocktype, coordinates are absolute
func Fill(p *mcwss.Player, pos mctype.Position, x1 int, y1 int, z1 int, x2 int, y2 int, z2 int, blocktype string) {
	p.Exec(fmt.Sprintf("fill %d %d %d %d %d %d %s", int(pos.X)+x1, int(pos.Y)+y1, int(pos.Z)+z1, int(pos.X)+x2, int(pos.Y)+y2, int(pos.Z)+z2, blocktype), nil)
}

// Summon will spawn a named entity at the coordinates relative to the position passed
func Summon(p *mcwss.Player, pos mctype.Position, x int, y int, z int, entity string, name string) {
	p.Exec(fmt.Sprintf("summon %s %s %d %d %d", entity, name, int(pos.X)+x, int(pos.Y)+y, int(pos.Z)+z), nil)
}

// Summonpos will spawn a named entity in a random area close to the position passed - UniqueID check will prevent spawning an entity more than once
func Summonpos(p *mcwss.Player, clientset *kubernetes.Clientset, pos mctype.Position, entity string, name string) {
	fmt.Printf("IN summon 1111\n")
	if !Contains(uniqueIDs, name) {
		fmt.Printf("IN summon 2222\n")
		uniqueIDs = append(uniqueIDs, name)
		p.Exec(fmt.Sprintf("summon %s %s %d %d %d", entity, name, int(pos.X-1.5+3*rand.Float64()), int(pos.Y)-5, int(pos.Z-1.5+3*rand.Float64())), nil)

		time.Sleep(100 * time.Millisecond)
		fmt.Printf("IN summon 3333\n")
	} else {
		fmt.Printf("Entity %s already exists\n", name)
		//ReconcileMCtoKubeMob(p, clientset, 12)
	}
}

// Testforentity will search for a named entity
func Testforentity(p *mcwss.Player, name string) bool {
	result := false
	go func() {
		p.Exec(fmt.Sprintf("testfor @e[name=%s]", name), func(response *command.LocalPlayerName) {
			if response.StatusCode == 0 {
				result = true
			}
		})
	}()

	time.Sleep(100 * time.Millisecond)
	return result
}

// Actionbar will display a message to the player
func Actionbar(player *mcwss.Player, message string) {
	player.Exec(fmt.Sprintf("title %s actionbar %s", player.Name(), message), nil)
}

// Get Current Player Position
func GetPlayerPosition(player *mcwss.Player) {

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
}

// Get Namspace Positions
func SetNamespacesPosition() {
	namespacesp = []mctype.Position{
		{X: initpos.X - 11, Y: initpos.Y + 5, Z: initpos.Z - 11},
		{X: initpos.X - 11, Y: initpos.Y + 5, Z: initpos.Z - 5},
		{X: initpos.X - 5, Y: initpos.Y + 5, Z: initpos.Z - 11},
		{X: initpos.X - 5, Y: initpos.Y + 5, Z: initpos.Z - 5},
	}
}
