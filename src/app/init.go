package main

import (
	"fmt"

	"github.com/sandertv/mcwss"
)

// MOTD will display the title and subtitle
func MOTD(player *mcwss.Player) {
	player.Exec(fmt.Sprintf("title %s title KubeCraftAdmin", player.Name()), nil)
	player.Exec(fmt.Sprintf("title %s subtitle The Adventurer's Admin Tool", player.Name()), nil)
}

// InitArea will spawn an initial playing area with 4 animal pens and a beacon at the center (currently hardcoded) and set the init position
func InitArea(p *mcwss.Player) {
	// Create animal pens
	fmt.Println("Creating Animal Pens ", initpos)

	Fill(p, initpos, -20, -2, -20, 20, 15, 20, "air")
	Fill(p, initpos, -15, -2, -15, 15, -1, 15, "stone 4")
	Fill(p, initpos, -1, -2, -1, 1, -2, 1, "glass")
	Fill(p, initpos, 0, -2, 0, 0, -2, 0, "beacon")
	Fill(p, initpos, -14, -1, -14, 14, -1, 14, "air")
	Fill(p, initpos, -14, -2, -14, -2, -2, -2, "grass")
	Fill(p, initpos, -14, -1, -14, -2, -1, -2, "fence")
	Fill(p, initpos, -13, -1, -13, -9, -1, -9, "air")
	Fill(p, initpos, -13, -1, -7, -9, -1, -3, "air")
	Fill(p, initpos, -7, -1, -13, -3, -1, -9, "air")
	Fill(p, initpos, -7, -1, -7, -3, -1, -3, "air")
}
