package main

func (effect *VFX) handleVFX(start bool) {
	// handle VFX logic
	if effect.name == FlameEffectName {
		if !start {
			effect.isActive = false
			effect.ticks = 0
		} else {
			effect.isActive = true
			// handle flame related logic
			if effect.ticks < uint16((cap(effect.effectImgs)*5))-1 {
				effect.ticks += 1
			} else {
				effect.ticks = 0
			}
		}
		effect.actingImg = effect.effectImgs[effect.ticks/5]
	}
}

func getPrioritizedTargets(g *Game, isVampireSatisfied bool) []*Player {
	// If my player has given the vampire food, then as a return
	// the vampire attacks the skeletons first,
	// then attacks the girl and my player.
	if isVampireSatisfied {
		targets := []*Player{
			g.SkeletonA,
			g.SkeletonB,
			g.girl.Player,
			g.myPlayer,
		}
		return targets
	}

	// Vampire has not been fed, as a revenge,
	// the vampire first attacks the girl and me, then goes to attack the skeletons.
	targets := []*Player{
		g.girl.Player,
		g.myPlayer,
		g.SkeletonA,
		g.SkeletonB,
	}
	return targets
}

func handleVampireMotion(g *Game) {
	vampire := g.vampire
	handleFoodEatLogic(g, vampire.Player)
	if !vampire.isFed {
		return
	}

	var target *Player
	for _, currentTarget := range getPrioritizedTargets(g, vampire.isSatisfied) {
		if currentTarget == g.girl.Player && g.girl.hiding {
			continue
		}
		if currentTarget.isActive && !currentTarget.escaped {
			target = currentTarget
			break
		}
	}

	if target == nil {
		g.flame.handleVFX(false)
		return
	}

	// Vampire approaches the target
	if !checkCollision(vampire.Sprite, target.Sprite, -1) {
		if vampire.x < target.x {
			g.playerMoves(vampire.Player, Right, &target.x, &target.y, nil)
		} else if vampire.x > target.x {
			g.playerMoves(vampire.Player, Left, &target.x, &target.y, nil)
		}

		if vampire.y < target.y {
			g.playerMoves(vampire.Player, Down, &target.x, &target.y, nil)
		} else if vampire.y > target.y {
			g.playerMoves(vampire.Player, Up, &target.x, &target.y, nil)
		}

		// Position the flame VFX around the vampire,
		// It should be at the bottom of the vampire.
		// with a little offset upwards.
		g.flame.y = vampire.y - (g.flame.height - vampire.height) + 3 // 3 here is just an offset
		// centered horizontally.
		g.flame.x = vampire.x - ((g.flame.width - vampire.width) / 2)

		target.handlePlayerDamage(false)
		g.flame.handleVFX(false)
	} else {
		// Vampire attacks the target
		g.flame.handleVFX(true)
		target.handlePlayerDamage(true)
	}

}
