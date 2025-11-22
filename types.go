package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const LayoutWidth float32 = 320  //400
const LayoutHeight float32 = 240 // 300
const TILE_SIZE int = 16

// Constant values related to motion:
const JumpPower float32 = 5.0
const Gravity float32 = 0.5
const AnimationTickDelay uint8 = 15

var DamagedRedColor color.RGBA = color.RGBA{255, 100, 100, 255}

type Resetter interface {
	Reset()
}

type SceneId uint8
type Direction uint8
type SpriteName string

const (
	Start SceneId = iota
	Playing
	Won
	Lost
)

const (
	Up Direction = iota
	Down
	Left
	Right
)

const (
	MyPlayerName    SpriteName = "inspector"
	GirlName        SpriteName = "girl"
	VampireName     SpriteName = "vampire"
	SkeletonA       SpriteName = "skeletonA"
	SkeletonB       SpriteName = "skeletonB"
	SkeletonC       SpriteName = "skeletonC"
	DoorA           SpriteName = "doorA"
	DoorB           SpriteName = "doorB"
	DoorC           SpriteName = "doorC"
	FlameEffectName SpriteName = "falmevfx"
	Blood           SpriteName = "blood"
	Meat            SpriteName = "meat"
	Shrimp          SpriteName = "sushi"
)

type Sprite struct {
	img           *ebiten.Image
	x, y          float32 // current sprite (x, y) position
	dX, dY        float32
	width, height float32
	actingImg     *ebiten.Image // current displayed img
	transparency  float32
	name          SpriteName
	isActive      bool
}

type Food struct {
	*Sprite
	eaten bool
}

type VFX struct {
	*Sprite
	effectImgs []*ebiten.Image
	ticks      uint16
}

type Door struct {
	*Sprite
	open bool
}

type Player struct {
	*Sprite
	hitTimer       uint16
	health         float32
	speed          float32
	animationTimer uint8
	animatedimgs   map[string]*ebiten.Image
	escaped        bool
}

type Vampire struct {
	*Player
	isSatisfied bool
	isFed       bool
	targets     []*Player // Players (in order) the vampire attacks
}

type Girl struct {
	*Player
	hiding bool
	door   *Door
}

type Game struct {
	currentScene  SceneId
	players       []*Player
	myPlayer      *Player
	girl          *Girl
	vampire       *Vampire
	SkeletonA     *Player
	SkeletonB     *Player
	doors         []*Door
	exit          *Sprite
	food          []*Food
	flame         *VFX
	showDialogFor SpriteName
	tilemapJSON   *TilemapJSON
	tilemapImg    *ebiten.Image
}
