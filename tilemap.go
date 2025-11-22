package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// data we want for one layer in our list of layers
type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

func DrawMap(screen *ebiten.Image, g *Game) {
	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.tilemapJSON.Layers {
		// loop over the tiles in the layer data
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}

			// get tile position then convert it to pixel position
			x := (index % layer.Width) * TILE_SIZE
			y := (index / layer.Width) * TILE_SIZE

			// Draw tile at x, y position.
			opts.GeoM.Translate(float64(x), float64(y))

			// map tile Id to its position in the actual png img.
			col := (id - 1) % (g.tilemapImg.Bounds().Dx() / TILE_SIZE)
			row := (id - 1) / (g.tilemapImg.Bounds().Dx() / TILE_SIZE)

			// draw the tile cropping out the tile that we want from the img.
			screen.DrawImage(
				getSubImage(g.tilemapImg, row, col, TILE_SIZE, TILE_SIZE),
				&opts,
			)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}
}

// opens the file, parses it, and returns the json object
func NewTilemapJSON(filepath string) *TilemapJSON {
	contents, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)

	if err != nil {
		log.Fatal(err)
	}

	return &tilemapJSON
}
