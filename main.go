package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func handleDoorOpen(g *Game) {
	if !g.girl.hiding {
		return
	}

	for _, door := range g.doors {
		// Check if player approaches one of the doors
		if checkCollision(g.myPlayer.Sprite, door.Sprite, 1) {
			// listen for the key stroke O (for open) so player can open that door.
			if inpututil.IsKeyJustPressed(ebiten.KeyO) {
				door.open = true
			}
		}

		if door.open {
			door.transparency = max(door.transparency-0.03, 0)
			// The girl should be free (not hiding)
			// if her door (i.e. the one she is hiding behind)
			// has opened completely i.e. its transparency=0
			if door.transparency == 0 {
				door.isActive = false
				if door == g.girl.door {
					fmt.Println("Player opened the RIGHT door !")
					g.girl.hiding = false
				}
			}
		}
	}
}

func getMyPlayerCollidors(g *Game) []*Sprite {
	colliders := []*Sprite{}
	for _, door := range g.doors {
		if door.isActive && !door.open {
			colliders = append(colliders, door.Sprite)
		}
	}
	return colliders
}

func handleFoodEatLogic(g *Game, potentialEater *Player) {
	for _, food := range g.food {
		if !food.isActive {
			continue
		}
		if food.eaten {
			food.transparency = max(food.transparency-0.03, 0)
			if food.transparency == 0 {
				food.isActive = false
			}
		} else {
			if checkCollision(potentialEater.Sprite, food.Sprite, 0) {
				// Vampire only eats once.
				if potentialEater == g.vampire.Player && !g.vampire.isFed {
					// Make vampire eat this food
					food.eaten = true
					fmt.Println("vampire ate", food.name)
					if food.name == Blood {
						fmt.Println("Vampire is now satisifed, since he ate ", food.name)
						g.vampire.isSatisfied = true
					}
					g.vampire.isFed = true
				}

				if potentialEater == g.myPlayer {
					if ebiten.IsKeyPressed(ebiten.KeyE) {
						food.eaten = true
						if food.name == Blood {
							// If my player eats blood,
							// decrement their health by 40%.
							g.myPlayer.health -= g.myPlayer.health * 0.4
						} else {
							// If it's other food
							// increment their health by 20%.
							g.myPlayer.health += g.myPlayer.health * 0.2
						}
					}
				}
				return
			}
		}

	}
}

func handleMyPlayerMotion(g *Game) {
	collidors := getMyPlayerCollidors(g)

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.playerMoves(g.myPlayer, Up, nil, nil, collidors)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyW) {
		g.myPlayer.playerStops("idle_back")
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.playerMoves(g.myPlayer, Down, nil, nil, collidors)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyDown) || inpututil.IsKeyJustReleased(ebiten.KeyS) {
		g.myPlayer.playerStops("idle_front")
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerMoves(g.myPlayer, Right, nil, nil, collidors)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyD) {
		g.myPlayer.playerStops("idle_right")
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerMoves(g.myPlayer, Left, nil, nil, collidors)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyA) {
		g.myPlayer.playerStops("idle_left")
	}
	// Check if my player has escaped
	g.myPlayer.hasPlayerEscaped(g)
}

func updateGameResult(g *Game) {
	if !g.myPlayer.isActive || !g.girl.isActive {
		// If my player or the girl dies,
		// You lose
		g.currentScene = Lost
	} else if g.myPlayer.escaped && g.girl.escaped {
		// If both my player and the girl manage to escape alive,
		// You win
		g.currentScene = Won
	}
}

func (g *Game) Update() error {
	if g.currentScene == Start {
		// start game
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.currentScene = Playing
		}
	}

	if g.currentScene == Playing {
		updateGameResult(g)
		handleMyPlayerMotion(g)
		handleFoodEatLogic(g, g.myPlayer)
		handleGirlMotion(g)
		handleSkeletonsMotion(g)
		handleVampireMotion(g)
		handleDoorOpen(g)
		handlePlayersDeath(g)
		handleDialogInfo(g)
	}

	if g.currentScene == Won || g.currentScene == Lost {
		// Play Again
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			*g = *NewGame(Playing)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	if g.currentScene == Start {
		// draw start screen, e.g. game intro & instructs
		drawStartScreen(screen)
	} else if g.currentScene == Won || g.currentScene == Lost {
		drawEndScreen(screen, g.currentScene == Won)
	} else {
		//DrawMap(screen, g)
		drawOpts := ebiten.DrawImageOptions{}

		drawActiveSprite(screen, g.exit, &drawOpts)
		drawActiveSprite(screen, g.flame.Sprite, &drawOpts)

		for _, food := range g.food {
			drawActiveSprite(screen, food.Sprite, &drawOpts)
		}

		// Sort Players in order of their Y position
		// This determines the draw order for each frame.
		sortPlayersByDrawOrder(g.players)
		for _, player := range g.players {
			// Flash red every 10 frames, for a player the gets attacked.
			if player.isActive && player.hitTimer > 0 && player.hitTimer%10 < 5 {
				drawOpts.ColorScale.ScaleWithColor(DamagedRedColor)
			}
			drawActiveSprite(screen, player.Sprite, &drawOpts)
		}

		for _, door := range g.doors {
			drawActiveSprite(screen, door.Sprite, &drawOpts)
		}

		if g.showDialogFor == VampireName {
			ebitenutil.DebugPrintAt(
				screen,
				"Vampire: I'm hungry !",
				3, 50,
			)
		}

		if g.showDialogFor == SkeletonA {
			ebitenutil.DebugPrintAt(
				screen,
				"Skeleton: Haha, don't worry, we won't harm you",
				3, int(LayoutHeight)-20,
			)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(LayoutWidth), int(LayoutHeight)
}

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("GoInspector - Created with enthusaism by Bedo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame(Start)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
