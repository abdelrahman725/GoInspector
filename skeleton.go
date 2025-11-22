package main

func handleSkeletonsMotion(g *Game) {
	// Skeletons should start moving once the girl is not hiding and hasn't escaped.
	if g.girl.hiding || g.girl.escaped {
		return
	}

	// Skeleton B blocks the exit.
	exitMiddle := (g.exit.x + (g.exit.width / 2)) - (g.SkeletonB.width / 2)
	exitTop := g.exit.y - g.exit.height - 4 // 4 is just an offset to leave space

	if g.SkeletonB.x < exitMiddle {
		g.playerMoves(g.SkeletonB, Right, &exitMiddle, &exitTop, nil)
	} else if g.SkeletonB.x > exitMiddle {
		g.playerMoves(g.SkeletonB, Left, &exitMiddle, &exitTop, nil)
	} else if g.SkeletonB.y < exitTop {
		g.playerMoves(g.SkeletonB, Down, &exitMiddle, &exitTop, nil)
	} else if g.SkeletonB.y > exitTop {
		g.playerMoves(g.SkeletonB, Up, &exitMiddle, &exitTop, nil)
	}

	girlCollidesWithSkeletonA := false
	girlCollidesWithSkeletonB := false

	if checkCollision(g.SkeletonB.Sprite, g.girl.Sprite, 2) {
		girlCollidesWithSkeletonB = true
	}

	// Skeleton A chases and attacks the girl.
	if checkCollision(g.SkeletonA.Sprite, g.girl.Sprite, -5) {
		girlCollidesWithSkeletonA = true
	} else {
		collidors := []*Sprite{g.exit}
		if g.SkeletonA.y < g.girl.y {
			g.playerMoves(g.SkeletonA, Down, &g.girl.x, &g.girl.y, collidors)
		} else if g.SkeletonA.y > g.girl.y {
			g.playerMoves(g.SkeletonA, Up, &g.girl.x, &g.girl.y, collidors)
		}
		if g.SkeletonA.x < g.girl.x {
			g.playerMoves(g.SkeletonA, Right, &g.girl.x, &g.girl.y, collidors)
		} else if g.SkeletonA.x > g.girl.x {
			g.playerMoves(g.SkeletonA, Left, &g.girl.x, &g.girl.y, collidors)
		}
	}

	if girlCollidesWithSkeletonA || girlCollidesWithSkeletonB {
		g.girl.handlePlayerDamage(true)
	} else {
		g.girl.handlePlayerDamage(false)
	}
}
