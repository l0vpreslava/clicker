package main

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bird struct {
	position          rl.Vector2
	velocity          rl.Vector2
	direction         int
	animationFrame    int
	animationTimer    float32
	animationInterval float32
}

type KilledBird struct {
	position  rl.Vector2
	timer     float32
	direction int
}

func (game *Game) updateState() {
	dt := rl.GetFrameTime()
	game.timer += dt

	if game.currentState == InGame {
		rl.HideCursor()
	} else {
		rl.ShowCursor()
	}

	switch game.currentState {
	case InGame:
		game.spawnBirds(dt)
		game.updateBirds(dt)
		game.updateAnimations(dt)
	}
}

func (game *Game) spawnBirds(dt float32) {
	if game.timer-game.lastSpawnTime >= game.spawnInterval {
		var x, y float32
		var velocity rl.Vector2
		var direction int
		y = rand.Float32()*(WorldHeight-2.0*BirdRadius) + BirdRadius

		switch rand.IntN(2) {
		case 0:
			x = -BirdRadius
			velocity = rl.NewVector2(1, 0)
			direction = Right

		case 1:
			x = float32(WorldWidth) + BirdRadius
			velocity = rl.NewVector2(-1, 0)
			direction = Left
		}

		speed := rand.Float32()*300 + 100
		velocity = rl.Vector2Scale(velocity, speed)

		game.birds = append(game.birds, Bird{
			position:          rl.NewVector2(x, y),
			velocity:          velocity,
			direction:         direction,
			animationFrame:    0,
			animationTimer:    0,
			animationInterval: 0.1,
		})

		game.lastSpawnTime = game.timer

		if game.spawnInterval > 0.5 {
			game.spawnInterval -= 0.1 * dt
		}
	}
}

func (game *Game) updateBirds(dt float32) {
	_ = dt

	var newBirds []Bird
	for _, bird := range game.birds {
		if (bird.direction == Right && bird.position.X > float32(WorldWidth)+BirdRadius) ||
			(bird.direction == Left && bird.position.X < 0-BirdRadius) {
			game.escapedBird++
			rl.PlaySound(game.assets.bannerSound)
			game.showBanner = !game.showBanner
			game.bannerTimer = BannerTimeMax
		} else {
			newBirds = append(newBirds, bird)
		}
	}

	game.birds = newBirds

	if game.escapedBird >= 10 {
		for _, bird := range game.birds {
			game.killedBirds = append(game.killedBirds, KilledBird{position: bird.position, timer: 0})
		}
		game.birds = make([]Bird, 0)
		game.currentState = GameOver
	}
}
