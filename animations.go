package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (game *Game) updateAnimations(dt float32) {
	game.updateAliveBirdsAnimations(dt)
	game.updateKilledBirdsAnimations(dt)
	game.updateBannerAnimations(dt)
}

func (game *Game) updateAliveBirdsAnimations(dt float32) {
	for i := range game.birds {
		game.birds[i].position = rl.Vector2Add(game.birds[i].position, rl.Vector2Scale(game.birds[i].velocity, dt))

		game.birds[i].animationTimer += dt
		if game.birds[i].animationTimer >= game.birds[i].animationInterval {
			game.birds[i].animationFrame = (game.birds[i].animationFrame + 1) % 2
			game.birds[i].animationTimer = 0
		}
	}
}

func (game *Game) updateKilledBirdsAnimations(dt float32) {
	for i := range game.killedBirds {
		const fallSpeed = -14.0    // u
		const acceleration = 100.0 // a

		bird := &game.killedBirds[i]

		bird.timer += dt // t
		// s = ut + 1/2*at^2
		bird.position.Y += fallSpeed*bird.timer + 0.5*acceleration*bird.timer*bird.timer
	}

	newKilledBirds := make([]KilledBird, 0)
	for _, killedBird := range game.killedBirds {
		if killedBird.position.Y < float32(WorldHeight) {
			newKilledBirds = append(newKilledBirds, killedBird)
		}
	}
	game.killedBirds = newKilledBirds
}

func (game *Game) updateBannerAnimations(dt float32) {
	if game.showBanner {
		game.bannerTimer -= dt
		if game.bannerTimer <= 0 {
			game.showBanner = false
		}
	}

	game.bannerAlpha = game.bannerTimer / BannerTimeMax
}
