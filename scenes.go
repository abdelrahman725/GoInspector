package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func drawStartScreen(screen *ebiten.Image) {
	logoImg, _ := customReadImage("assets/images/logo.png")

	// Calculate scaling for logo to draw it at the top of the screen.
	logoWidth := logoImg.Bounds().Dx()
	logoHeight := logoImg.Bounds().Dy()
	logoX := float64(LayoutWidth) / float64(logoWidth)
	// Logo takes 20% of screen height.
	logoY := float64(LayoutHeight*0.20) / float64(logoHeight)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(logoX, logoY)
	screen.DrawImage(logoImg, opts)
	opts.GeoM.Reset()

	introText := "The inspector and the hidden girl must escape safely\n" +
		"through the exit. \n" +
		"How to do that ? Well, that's the game !\n\n" +
		"------------------------------------------\n" +
		"Use arrows or the keys W, A, S, D to move.\n" +
		"Press O to open a door.\n" +
		"Press E to eat a food.\n" +
		"You can also push the food.\n" +
		"------------------------------------------\n" +
		"Press Enter to start..."

	ebitenutil.DebugPrintAt(screen, introText, 3, 50)
}

func drawEndScreen(screen *ebiten.Image, wonGame bool) {
	if wonGame {
		ebitenutil.DebugPrintAt(
			screen,
			"Congratulations !\nYou are a true savior and also smart !\n\nPress Enter to Play again",
			20, 20,
		)
	} else {
		ebitenutil.DebugPrintAt(
			screen,
			"Opps Try again!\n\nPress Enter to Play again",
			20, 20,
		)
	}
}
