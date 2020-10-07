package main

import (
	"fmt"
	"math"
	"time"

	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/mctype"
	"github.com/sandertv/mcwss/protocol/command"
)

func agentFollow(a *mcwss.Agent, p *mcwss.Player) {

	jump := 0
	fx := 0
	fz := 0

	for { // Action loop for Agent

		time.Sleep(500 * time.Millisecond)
		a.Position(func(aPos mctype.Position) {
			p.Exec(fmt.Sprintf("testforblock %d %d %d air", int(aPos.X), int(aPos.Y)-1, int(aPos.Z)), func(response *command.LocalPlayerName) {
				if (jump == 0) && (response.StatusCode == 0) {
					a.Move(mctype.Down, 1)
					fmt.Println("Fall!")
				} else {
					p.Exec(fmt.Sprintf("testforblock %d %d %d air", int(aPos.X)+fx, int(aPos.Y), int(aPos.Z)+fz), func(response *command.LocalPlayerName) {
						if response.StatusCode != 0 {
							a.Move(mctype.Up, 1)
							jump = 1
							fmt.Printf("%d %d", fx, fz)
							fmt.Println("Jump!")
						} else {
							jump = 0
							a.Rotation(func(aRot float64) {
								p.Position(func(pPos mctype.Position) {
									// Find out desired orientation to face the player
									var dRot float64
									xdist := pPos.X - aPos.X
									zdist := pPos.Z - aPos.Z
									if math.Abs(xdist) > math.Abs(zdist) {
										if xdist > 0 {
											dRot = -90
										} else {
											dRot = 90
										}
									} else {
										if zdist > 0 {
											dRot = 0
										} else {
											dRot = -180
										}
									}

									if dRot != aRot {
										fx = 0
										fz = 0
										// Decide to turn right or left
										if (dRot == 90) && (aRot == -180) {
											a.TurnLeft()
										} else if (dRot == -180) && (aRot == 90) {
											a.TurnRight()
										} else {
											if dRot > aRot {
												a.TurnRight()
											} else {
												a.TurnLeft()
											}
										}
									} else {
										if aRot == 90 {
											fx = -1
										}
										if aRot == -90 {
											fx = 1
										}
										if aRot == 0 {
											fz = 1
										}
										if aRot == -180 {
											fz = -1
										}
										a.Move(mctype.Forward, 1)
									}
								})

							})
						}
					})
				}
			})
		})
	}
}
