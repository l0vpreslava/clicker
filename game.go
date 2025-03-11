package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Menu = iota
	InSettings
	InGame
	Pause
	GameOver
)

const (
	Left  = -1
	Right = 1
)

const WorldWidth = 640
const WorldHeight = 360
const BannerTimeMax float32 = 0.5
const BirdRadius float32 = 27

type Game struct {
	ballColor     rl.Color
	birds         []Bird
	timer         float32
	spawnInterval float32
	lastSpawnTime float32
	killedBirds   []KilledBird
	score         int
	scores        []int
	escapedBird   int
	currentState  int
	settings      Settings
	assets        Assets
	showBanner    bool
	bannerTimer   float32
	shouldClose   bool
	bannerAlpha   float32
	TextColor     rl.Color
	RenderTexture rl.RenderTexture2D
}

func NewGame() Game {
	return Game{
		ballColor:     rl.GetColor(0xf96e61ff),
		birds:         make([]Bird, 0),
		timer:         0,
		spawnInterval: 1.5,
		lastSpawnTime: 0,
		killedBirds:   make([]KilledBird, 0),
		score:         0,
		scores:        make([]int, 0),
		currentState:  Menu,
		settings:      NewSettigs(),
		assets:        LoadAssets(),
		showBanner:    false,
		bannerTimer:   0,
		shouldClose:   false,
		bannerAlpha:   0,
		TextColor:     rl.GetColor(0xc655f6ff),
		RenderTexture: rl.LoadRenderTexture(WorldWidth, WorldHeight),
	}
}

func (game *Game) Destroy() {
	game.assets.Unload()
	rl.UnloadRenderTexture(game.RenderTexture)
}

func (game *Game) Update() {
	game.handleInput()
	game.updateState()

	rl.BeginDrawing()

	game.draw()
	game.handleUI()

	rl.EndDrawing()
}
