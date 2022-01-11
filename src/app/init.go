package main

import (
	"fmt"
	"github.com/sandertv/mcwss/mctype"
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
	p.Position(func(pos mctype.Position){
		fmt.Println("Creating Animal Pens ", pos)

		Fill(p, pos, -20, -2, -20, 20, 15, 20, "air")
		Fill(p, pos, -15, -2, -15, 15, -1, 15, "stone 4")
		Fill(p, pos, -1, -2, -1, 1, -2, 1, "glass")
		Fill(p, pos, 0, -2, 0, 0, -2, 0, "beacon")
		Fill(p, pos, -14, -1, -14, 14, -1, 14, "air")
		Fill(p, pos, -14, -2, -14, -2, -2, -2, "grass")
		Fill(p, pos, -14, -1, -14, -2, -1, -2, "fence")
		Fill(p, pos, -13, -1, -13, -9, -1, -9, "air")
		Fill(p, pos, -13, -1, -7, -9, -1, -3, "air")
		Fill(p, pos, -7, -1, -13, -3, -1, -9, "air")
		Fill(p, pos, -7, -1, -7, -3, -1, -3, "air")	
	})
}
