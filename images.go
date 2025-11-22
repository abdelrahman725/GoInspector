package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func drawActiveSprite(screen *ebiten.Image, sprite *Sprite, opts *ebiten.DrawImageOptions) {
	// Draw only active (i.e. currently existing in game) sprite.
	if sprite.isActive {
		opts.GeoM.Translate(float64(sprite.x), float64(sprite.y))
		opts.ColorScale.ScaleAlpha(sprite.transparency)
		screen.DrawImage(sprite.actingImg, opts)
	}
	opts.ColorScale.Reset()
	opts.GeoM.Reset()
}

func customReadImage(filepath string) (*ebiten.Image, image.Image) {
	imgPointer, imgData, err := ebitenutil.NewImageFromFile(filepath)

	if err != nil {
		log.Fatal(err)
	}
	return imgPointer, imgData
}

func getSubImage(img *ebiten.Image, row, col int, x_length, y_length int) *ebiten.Image {
	x_start := col * x_length
	y_start := row * y_length
	sub_img := img.SubImage(
		image.Rect(x_start, y_start, x_start+x_length, y_start+y_length),
	).(*ebiten.Image)

	return sub_img
}

func copyImage(original *ebiten.Image) *ebiten.Image {
	// Create a new, independent copy of an ebiten.Image.

	// 1. Create a new, empty image with the same dimensions as the original.
	newImage := ebiten.NewImage(original.Bounds().Dx(), original.Bounds().Dy())

	// 2. Draw the original image onto the new image.
	// This performs a pixel-by-pixel copy.
	newImage.DrawImage(original, &ebiten.DrawImageOptions{})

	return newImage
}

func constructVfxImgs(flame *VFX) {
	// FlameSpriteSheet dimensions:
	// width  = 200, 8 images, so each one is 25 px width.
	// height = 30
	for i := range cap(flame.effectImgs) {
		flame.effectImgs[i] = getSubImage(flame.img, 0, i, 25, 30)
	}
}

func constructAnimationImgs(g *Game) {
	for _, player := range g.players {
		if player.name == SkeletonB {
			continue
		}
		player.animatedimgs["idle_front"] = getSubImage(player.img, 0, 0, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["idle_back"] = getSubImage(player.img, 0, 1, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["idle_left"] = getSubImage(player.img, 0, 2, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["idle_right"] = getSubImage(player.img, 0, 3, TILE_SIZE, TILE_SIZE)

		player.animatedimgs["running_front_1"] = getSubImage(player.img, 1, 0, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["running_front_2"] = getSubImage(player.img, 3, 0, TILE_SIZE, TILE_SIZE)

		player.animatedimgs["running_back_1"] = getSubImage(player.img, 1, 1, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["running_back_2"] = getSubImage(player.img, 3, 1, TILE_SIZE, TILE_SIZE)

		player.animatedimgs["running_left_1"] = getSubImage(player.img, 1, 2, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["running_left_2"] = getSubImage(player.img, 3, 2, TILE_SIZE, TILE_SIZE)

		player.animatedimgs["running_right_1"] = getSubImage(player.img, 1, 3, TILE_SIZE, TILE_SIZE)
		player.animatedimgs["running_right_2"] = getSubImage(player.img, 3, 3, TILE_SIZE, TILE_SIZE)
		player.actingImg = player.animatedimgs["idle_front"]

	}
	g.SkeletonB.animatedimgs = g.SkeletonA.animatedimgs
	g.SkeletonB.actingImg = g.SkeletonB.animatedimgs["idle_front"]
}

func createExitSprite(width, height int) *ebiten.Image {
	// Create width x height image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Colors - classic exit sign
	bgColor := color.RGBA{0, 128, 0, 255}         // Green background
	borderColor := color.RGBA{255, 255, 255, 255} // White border
	textColor := color.RGBA{255, 255, 255, 255}   // White text

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw border (1px white)
	for x := 0; x < width; x++ {
		img.Set(x, 0, borderColor)        // Top
		img.Set(x, height-1, borderColor) // Bottom
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, borderColor)       // Left
		img.Set(width-1, y, borderColor) // Right
	}

	// Draw "EXIT" text pixel by pixel (simple 3x5 font)
	// Total width: 4 letters * 3px + 3 spaces * 1px = 15px
	// Start at x = (32 - 15) / 2 = 8.5, rounded to 8
	startX := 8
	y := 6

	drawPixelLetter(img, 'E', startX, y, textColor)    // x=8
	drawPixelLetter(img, 'X', startX+5, y, textColor)  // x=12
	drawPixelLetter(img, 'I', startX+10, y, textColor) // x=16
	drawPixelLetter(img, 'T', startX+15, y, textColor) // x=20

	// Convert to Ebitengine image
	return ebiten.NewImageFromImage(img)
}

func drawPixelLetter(img *image.RGBA, letter rune, x, y int, c color.RGBA) {
	// Helper function to draw pixel letters using a single function
	switch letter {
	case 'E':
		// E shape (3x5)
		img.Set(x, y, c) // Top row
		img.Set(x+1, y, c)
		img.Set(x+2, y, c)
		img.Set(x, y+1, c) // Left side
		img.Set(x, y+2, c)
		img.Set(x+1, y+2, c) // Middle row
		img.Set(x, y+3, c)   // Left side
		img.Set(x, y+4, c)
		img.Set(x+1, y+4, c) // Bottom row
		img.Set(x+2, y+4, c)

	case 'X':
		// X shape (3x5)
		img.Set(x, y, c) // Top corners
		img.Set(x+2, y, c)
		img.Set(x+1, y+1, c) // Center diagonal
		img.Set(x+1, y+2, c)
		img.Set(x+1, y+3, c)
		img.Set(x, y+4, c) // Bottom corners
		img.Set(x+2, y+4, c)

	case 'I':
		// I shape (3x5)
		img.Set(x, y, c) // Top row
		img.Set(x+1, y, c)
		img.Set(x+2, y, c)
		img.Set(x+1, y+1, c) // Center column
		img.Set(x+1, y+2, c)
		img.Set(x+1, y+3, c)
		img.Set(x, y+4, c) // Bottom row
		img.Set(x+1, y+4, c)
		img.Set(x+2, y+4, c)

	case 'T':
		// T shape (3x5)
		img.Set(x, y, c) // Top row
		img.Set(x+1, y, c)
		img.Set(x+2, y, c)
		img.Set(x+1, y+1, c) // Center column
		img.Set(x+1, y+2, c)
		img.Set(x+1, y+3, c)
		img.Set(x+1, y+4, c)
	}
}
