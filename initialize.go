package main

import "github.com/hajimehoshi/ebiten/v2"

// var tilemapJSON = NewTilemapJSON("assets/maps/tiles.json")
// var tilemapImg, _ = customReadImage("assets/images/Pipes.png")
var myPlayerImg, _ = customReadImage("assets/images/characters/InspectorSpriteSheet.png")
var vampireImg, _ = customReadImage("assets/images/characters/VampireSpriteSheet.png")
var girlImg, _ = customReadImage("assets/images/characters/PrincessSpriteSheet.png")
var SkeletonImg, _ = customReadImage("assets/images/characters/SkeletonSpriteSheet.png")
var doorImg, _ = customReadImage("assets/images/Door.png")
var flameImg, _ = customReadImage("assets/images/FlameSpriteSheet.png")
var exitSprite = createExitSprite(32, 17)
var bloodImg, _ = customReadImage("assets/images/food/Blood.png")
var meatImg, _ = customReadImage("assets/images/food/Meat.png")
var shrimpImg, _ = customReadImage("assets/images/food/Shrimp.png")

func NewGame(scene SceneId) *Game {
	// Initialize our non-player spirtes
	meat := Food{
		Sprite: &Sprite{
			actingImg:    meatImg,
			transparency: 1,
			x:            100,
			y:            150,
			width:        16,
			height:       16,
			isActive:     true,
			name:         Meat,
		},
	}

	blood := Food{
		Sprite: &Sprite{
			actingImg:    bloodImg,
			transparency: 1,
			x:            140,
			y:            150,
			width:        16,
			height:       16,
			isActive:     true,
			name:         Blood,
		},
	}

	shrimp := Food{
		Sprite: &Sprite{
			actingImg:    shrimpImg,
			transparency: 1,
			x:            180,
			y:            150,
			width:        16,
			height:       16,
			isActive:     true,
			name:         Shrimp,
		},
	}

	food := []*Food{&meat, &blood, &shrimp}

	flameVfx := VFX{
		Sprite: &Sprite{
			img:          flameImg,
			transparency: 1,
			name:         FlameEffectName,
			width:        25,
			height:       30,
			isActive:     false,
		},
		effectImgs: make([]*ebiten.Image, 8),
	}

	exit := Sprite{
		actingImg: exitSprite,
		// place the exit in the bottom right corner.
		x:            LayoutWidth - 32,
		y:            LayoutHeight - 17,
		width:        32,
		height:       17,
		transparency: 1,
		isActive:     true,
	}

	doorA := Door{
		Sprite: &Sprite{
			actingImg:    doorImg,
			width:        32,
			height:       32,
			x:            100,
			y:            15,
			transparency: 1.0,
			isActive:     true,
			name:         DoorA,
		},
	}

	doorB := Door{
		Sprite: &Sprite{
			actingImg:    copyImage(doorImg),
			width:        32,
			height:       32,
			x:            152,
			y:            15,
			transparency: 1.0,
			isActive:     true,
			name:         DoorB,
		},
	}

	doorC := Door{
		Sprite: &Sprite{
			actingImg:    copyImage(doorImg),
			width:        32,
			height:       32,
			x:            204,
			y:            15,
			transparency: 1.0,
			isActive:     true,
			name:         DoorC,
		},
	}

	doors := []*Door{&doorA, &doorB, &doorC}

	// Initialize our players
	myPlayer := Player{
		Sprite: &Sprite{
			img:          myPlayerImg,
			x:            100,
			y:            100,
			width:        float32(TILE_SIZE),
			height:       float32(TILE_SIZE),
			name:         MyPlayerName,
			transparency: 1.0,
			isActive:     true,
		},
		speed:        1.5,
		animatedimgs: make(map[string]*ebiten.Image),
		// my player dies after 180 ticks (3 seconds) if he gets attacked
		// Note these ticks/seconds are not necessarily continuous,
		// but the total acquired time he is exposed to damage.
		health: 180,
	}

	vampire := Vampire{
		Player: &Player{
			Sprite: &Sprite{
				img:          vampireImg,
				x:            7,
				y:            7,
				width:        float32(TILE_SIZE),
				height:       float32(TILE_SIZE),
				name:         VampireName,
				transparency: 1.0,
				isActive:     true,
			},
			speed:        1.3,
			animatedimgs: make(map[string]*ebiten.Image),
			health:       300,
		},
		targets: make([]*Player, 5),
	}

	girl := Girl{
		Player: &Player{
			Sprite: &Sprite{
				img:          girlImg,
				width:        float32(TILE_SIZE),
				height:       float32(TILE_SIZE),
				name:         GirlName,
				transparency: 1.0,
				isActive:     true,
			},
			speed:        0.8,
			animatedimgs: make(map[string]*ebiten.Image),
			health:       60, // girl dies after 50 ticks (i.e. accumulated time of 1 second) if she gets attacked.
		},
		hiding: true,
	}

	skeletonPlayerA := Player{
		Sprite: &Sprite{
			img:          SkeletonImg,
			x:            10,
			y:            200,
			width:        float32(TILE_SIZE),
			height:       float32(TILE_SIZE),
			name:         SkeletonA,
			transparency: 1.0,
			isActive:     true,
		},
		speed:        1,
		animatedimgs: make(map[string]*ebiten.Image),
		health:       100,
	}

	skeletonPlayerB := Player{
		Sprite: &Sprite{
			img:          SkeletonImg,
			x:            40,
			y:            200,
			width:        float32(TILE_SIZE),
			height:       float32(TILE_SIZE),
			name:         SkeletonB,
			transparency: 1.0,
			isActive:     true,
		},
		speed:        1,
		animatedimgs: make(map[string]*ebiten.Image),
		health:       100,
	}

	game := &Game{
		currentScene: scene,
		players: []*Player{
			&myPlayer,
			vampire.Player,
			girl.Player,
			&skeletonPlayerA,
			&skeletonPlayerB,
		},
		myPlayer:  &myPlayer,
		girl:      &girl,
		vampire:   &vampire,
		SkeletonA: &skeletonPlayerA,
		SkeletonB: &skeletonPlayerB,
		doors:     doors,
		flame:     &flameVfx,
		exit:      &exit,
		food:      food,
		//tilemapJSON:   tilemapJSON,
		//tilemapImg:    tilemapImg,
		showDialogFor: "",
	}

	constructAnimationImgs(game)
	constructVfxImgs(&flameVfx)

	girl.placeBehindRandomDoor(doors)
	return game
}
