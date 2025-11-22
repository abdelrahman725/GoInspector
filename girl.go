package main

import "math/rand"

func handleGirlMotion(g *Game) {
	if !g.girl.isActive || g.girl.hiding || g.girl.escaped {
		return
	}
	// The girl is free, now she must escape by running towards the exit.
	exitMiddle := (g.exit.x + (g.exit.width / 2)) - (g.girl.width / 2)
	if g.girl.x < exitMiddle {
		g.playerMoves(g.girl.Player, Right, &exitMiddle, &g.exit.y, nil)
	} else if g.girl.x > exitMiddle {
		g.playerMoves(g.girl.Player, Left, &exitMiddle, &g.exit.y, nil)
	} else if g.girl.y < g.exit.y {
		g.playerMoves(g.girl.Player, Down, &exitMiddle, &g.exit.y, nil)
	} else if g.girl.y > g.exit.y {
		g.playerMoves(g.girl.Player, Up, &exitMiddle, &g.exit.y, nil)
	}

	// Check if the girl has escaped
	g.girl.hasPlayerEscaped(g)
}

func (girl *Girl) placeBehindRandomDoor(doors []*Door) {
	// First choose a random door.
	randomIndex := rand.Intn(len(doors))
	girl.door = doors[randomIndex]

	// Place the girl behind (already done by drawing order) that random door
	// in a vertically and horizontally centered position, relative to the door position.
	girl.x = doors[randomIndex].x + ((doors[randomIndex].width - girl.width) / 2)
	girl.y = doors[randomIndex].y + ((doors[randomIndex].height - girl.height) / 2)
}
