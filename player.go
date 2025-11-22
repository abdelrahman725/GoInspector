package main

import (
	"sort"
)

func sortPlayersByDrawOrder(players []*Player) {
	// Sprites that are lower on the screen (larger y) are drawn on top
	// and Sprites that are higher on the screen (smaller y) are drawn on bottom
	sort.Slice(players, func(i, j int) bool {
		return players[i].y < players[j].y
	})
}

func (g *Game) handlePushingLogicOnlyforMyPlayer(player *Player) bool {
	// Hanldes Pushing system logic.
	// Returns bool indicating whether player succeeded to push
	// a sprite or not.

	if player != g.myPlayer {
		return true
	}

	futuerPlayerSprite := &Sprite{
		x:      player.x + player.dX,
		y:      player.y + player.dY,
		width:  player.width,
		height: player.height,
	}

	// Loop over pushable sprites (only food for now)
	for _, food := range g.food {
		if !food.isActive || food.eaten {
			continue
		}
		if checkCollision(food.Sprite, futuerPlayerSprite, -4) {
			if g.tryPushSprite(food.Sprite, player.dX, player.dY) {
				food.x += player.dX
				food.y += player.dY
				return true
			}
			return false
		}
	}
	return true

}

func (g *Game) tryPushSprite(sprite *Sprite, dx, dy float32) bool {
	// Future movable sprite position
	newX := sprite.x + dx
	newY := sprite.y + dy

	// Check collision with edges.
	if newX < 0 || newY < 0 || newX+sprite.width > LayoutWidth || newY+sprite.height > LayoutHeight {
		return false
	}

	futuerSprite := &Sprite{
		x:      newX,
		y:      newY,
		width:  sprite.width,
		height: sprite.height,
	}

	// check collision with food items
	for _, food := range g.food {
		if food.isActive && food.Sprite != sprite && checkCollision(food.Sprite, futuerSprite, 0) {
			return false
		}
	}

	// check collision with active doors.
	for _, door := range g.doors {
		if door.isActive && checkCollision(door.Sprite, futuerSprite, 0) {
			return false
		}
	}

	return true
}

func (player *Player) playerStops(idle_img string) {
	if idle_img != "" {
		player.actingImg = player.animatedimgs[idle_img]
	}
	player.animationTimer = 0
	player.dX = 0
	player.dY = 0
}

func (g *Game) playerMoves(
	player *Player,
	direction Direction,
	targetX *float32,
	targetY *float32,
	colliders []*Sprite,
) {
	toggleSameDirectionAnimation := func(animated_img_1, animated_img_2 string) {
		player.animationTimer += 1
		if player.animationTimer > AnimationTickDelay {
			player.actingImg = player.animatedimgs[animated_img_1]
		} else {
			player.actingImg = player.animatedimgs[animated_img_2]
		}
		if player.animationTimer > AnimationTickDelay*2 {
			player.animationTimer = 0
		}
	}

	if player.health == 0 || !player.isActive {
		return
	}

	if direction == Up {
		player.dX = 0
		if targetY != nil && player.y-player.speed < *targetY {
			player.dY = -Abs32(*targetY - player.y)
		} else {
			player.dY = -player.speed
		}
		toggleSameDirectionAnimation("running_back_1", "running_back_2")
	} else if direction == Down {
		player.dX = 0
		if targetY != nil && player.y+player.speed > *targetY {
			player.dY = Abs32(*targetY - player.y)
		} else {
			player.dY = player.speed
		}
		toggleSameDirectionAnimation("running_front_1", "running_front_2")

	} else if direction == Right {
		player.dY = 0
		if targetX != nil && player.x+player.speed > *targetX {
			player.dX = Abs32(*targetX - player.x)
		} else {
			player.dX = player.speed
		}
		toggleSameDirectionAnimation("running_right_1", "running_right_2")

	} else if direction == Left {
		player.dY = 0
		if targetX != nil && player.x-player.speed < *targetX {
			player.dX = -Abs32(*targetX - player.x)
		} else {
			player.dX = -player.speed
		}
		toggleSameDirectionAnimation("running_left_1", "running_left_2")
	}

	newY := player.y + player.dY
	newX := player.x + player.dX

	if newY <= LayoutHeight-player.height && newY > 1 {
		if g.handlePushingLogicOnlyforMyPlayer(player) {
			player.y = newY
		}
	}

	if newX < LayoutWidth-player.height && newX > 1 {
		if g.handlePushingLogicOnlyforMyPlayer(player) {
			player.x = newX
		}
	}

	preventCollision(player.Sprite, colliders, direction == Left || direction == Right)
}

func (player *Player) hasPlayerEscaped(g *Game) {
	// Checks if a player has escaped successfully by reaching the exit.
	// Note: Only called on the girl and my player

	// Since exit is positioned in the right bottom corner
	// It's enough to check that the player has passed through
	// the exit's upper and left side.
	if player.y >= g.exit.y && player.x >= g.exit.x {
		player.hitTimer = 0
		player.escaped = true
		//fmt.Println(player.name, "has escapped !")
	} else {
		player.escaped = false
		//fmt.Println(player.name, "NOT escaped!")
	}
}

func (target *Player) handlePlayerDamage(damage bool) {
	if !damage {
		target.hitTimer = 0
		return
	}

	if target.hitTimer > 0 {
		target.hitTimer -= 1
	} else {
		target.hitTimer = 30
	}

	if target.health > 0 {
		target.health -= 1
	}
}

func handlePlayersDeath(g *Game) {
	for _, attackedPlayer := range g.players {
		if attackedPlayer.health != 0 {
			continue
		}

		attackedPlayer.transparency = max(attackedPlayer.transparency-0.02, 0)
		// if player has faded completely, we can now mark them as dead.
		if attackedPlayer.transparency == 0 {
			attackedPlayer.isActive = false
		}
	}
}
