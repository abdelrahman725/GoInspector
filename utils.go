package main

import (
	"image"
)

func Abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func checkCollision(sprite_a *Sprite, sprite_b *Sprite, gap int) bool {
	spriteARect := image.Rect(
		int(sprite_a.x)-gap,
		int(sprite_a.y)-gap,
		int(sprite_a.x+sprite_a.width)+gap,
		int(sprite_a.y+sprite_a.height)+gap,
	)
	spriteBRect := image.Rect(
		int(sprite_b.x),
		int(sprite_b.y),
		int(sprite_b.x+sprite_b.width),
		int(sprite_b.y+sprite_b.height),
	)
	return spriteARect.Overlaps(spriteBRect)
}

func preventCollision(sprite *Sprite, colliders []*Sprite, horizontal bool) {
	if colliders == nil {
		return
	}

	spriteRect := image.Rect(
		int(sprite.x),
		int(sprite.y),
		int(sprite.x+sprite.width),
		int(sprite.y+sprite.height),
	)

	for _, collidor := range colliders {
		if sprite == collidor {
			continue
		}
		collidorRect := image.Rect(
			int(collidor.x),
			int(collidor.y),
			int(collidor.x+collidor.width),
			int(collidor.y+collidor.height),
		)
		if horizontal {
			if collidorRect.Overlaps(spriteRect) {
				if sprite.dX > 0 {
					sprite.x = float32(collidorRect.Min.X) - float32(TILE_SIZE)
				} else if sprite.dX < 0 {
					sprite.x = float32(collidorRect.Max.X)
				}
			}
		} else {
			if collidorRect.Overlaps(spriteRect) {
				if sprite.dY > 0 {
					sprite.y = float32(collidorRect.Min.Y) - float32(TILE_SIZE)
				} else if sprite.dY < 0 {
					sprite.y = float32(collidorRect.Max.Y)
				}
			}
		}
	}
}

func handleDialogInfo(g *Game) {
	vampire := g.vampire
	inspector := g.myPlayer

	if !vampire.isFed && checkCollision(vampire.Sprite, inspector.Sprite, 2) {
		g.showDialogFor = VampireName

	} else if g.girl.hiding && (checkCollision(g.SkeletonA.Sprite, inspector.Sprite, 2) || checkCollision(g.SkeletonB.Sprite, inspector.Sprite, 2)) {
		g.showDialogFor = SkeletonA
	} else {
		g.showDialogFor = ""
	}
}
