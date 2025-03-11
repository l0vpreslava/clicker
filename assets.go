package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Assets struct {
	font          rl.Font
	bgMenu        rl.Texture2D
	bgGame        rl.Texture2D
	banner        rl.Texture2D
	birdWingsUp   rl.Texture2D
	birdWingsDown rl.Texture2D
	birdDead      rl.Texture2D
	sight         rl.Texture2D
	dislike       rl.Texture2D
	selectButton  rl.Sound
	startGame     rl.Sound
	pause         rl.Sound
	bannerSound   rl.Sound
	shoot_1       rl.Sound
	shoot_2       rl.Sound
	birdDied_1    rl.Sound
	birdDied_2    rl.Sound
}

func LoadAssets() Assets {
	return Assets{
		font:          rl.LoadFontEx("assets/fonts/pixeleum-48.ttf", 96, nil, 0),
		bgMenu:        rl.LoadTexture("assets/sprites/bg-menu.png"),
		bgGame:        rl.LoadTexture("assets/sprites/bg-game.png"),
		banner:        rl.LoadTexture("assets/sprites/bg-rep.png"),
		dislike:       rl.LoadTexture("assets/sprites/hand.png"),
		birdWingsUp:   rl.LoadTexture("assets/sprites/bird1.png"),
		birdWingsDown: rl.LoadTexture("assets/sprites/bird2.png"),
		birdDead:      rl.LoadTexture("assets/sprites/bird3.png"),
		sight:         rl.LoadTexture("assets/sprites/sight.png"),
		bannerSound:   rl.LoadSound("assets/sounds/banner.wav"),
		birdDied_1:    rl.LoadSound("assets/sounds/bird-died-1.wav"),
		birdDied_2:    rl.LoadSound("assets/sounds/bird-died-2.wav"),
		pause:         rl.LoadSound("assets/sounds/pause.wav"),
		selectButton:  rl.LoadSound("assets/sounds/select-button.wav"),
		shoot_1:       rl.LoadSound("assets/sounds/shoot-1.wav"),
		shoot_2:       rl.LoadSound("assets/sounds/shoot-2.wav"),
		startGame:     rl.LoadSound("assets/sounds/start-game.wav"),
	}
}

func (assets *Assets) Unload() {
	rl.UnloadFont(assets.font)
	rl.UnloadTexture(assets.bgMenu)
	rl.UnloadTexture(assets.bgGame)
	rl.UnloadTexture(assets.banner)
	rl.UnloadTexture(assets.dislike)
	rl.UnloadTexture(assets.birdWingsUp)
	rl.UnloadTexture(assets.birdWingsDown)
	rl.UnloadTexture(assets.birdDead)
	rl.UnloadTexture(assets.sight)
	rl.UnloadSound(assets.selectButton)
	rl.UnloadSound(assets.startGame)
	rl.UnloadSound(assets.pause)
	rl.UnloadSound(assets.bannerSound)
	rl.UnloadSound(assets.shoot_1)
	rl.UnloadSound(assets.shoot_2)
	rl.UnloadSound(assets.birdDied_1)
	rl.UnloadSound(assets.birdDied_2)
}
