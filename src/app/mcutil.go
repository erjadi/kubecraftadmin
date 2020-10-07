package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/mctype"
	"github.com/sandertv/mcwss/protocol/command"
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

// Summonpos will spawn a named entity in a random area close to the position passed
func Summonpos(p *mcwss.Player, pos mctype.Position, entity string, name string) {
	if !Contains(uniqueIDs, name) {
		uniqueIDs = append(uniqueIDs, name)
		p.Exec(fmt.Sprintf("summon %s %s %d %d %d", entity, name, int(pos.X-1.5+3*rand.Float64()), int(pos.Y)-5, int(pos.Z-1.5+3*rand.Float64())), nil)
		time.Sleep(100 * time.Millisecond)
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

// Actionbar will display a message to the playerout

func Actionbar(player *mcwss.Player, message string) {
	player.Exec(fmt.Sprintf("title %s actionbar %s", player.Name(), message), nil)
}
