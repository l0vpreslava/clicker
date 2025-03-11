package main

import (
	"log"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (game *Game) handleInput() {
	if rl.WindowShouldClose() {
		game.shouldClose = true
	}

	if game.currentState == InGame {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			var sound rl.Sound
			if rand.IntN(2) == 0 {
				sound = game.assets.shoot_1
			} else {
				sound = game.assets.shoot_2
			}
			rl.PlaySound(sound)

			birds := make([]Bird, len(game.birds))
			index := 0
			for _, bird := range game.birds {
				if !rl.CheckCollisionPointCircle(getMouseWorldPos(), bird.position, BirdRadius) {
					birds[index] = bird
					index++
				} else {
					game.killedBirds = append(game.killedBirds, KilledBird{position: bird.position, timer: 0, direction: bird.direction})
					game.score++

					var sound rl.Sound
					if rand.IntN(2) == 0 {
						sound = game.assets.birdDied_1
					} else {
						sound = game.assets.birdDied_2
					}
					rl.PlaySound(sound)
				}
			}

			game.birds = birds[:index]
		}
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		rl.PlaySound(game.assets.pause)
		if game.currentState == InGame {
			game.currentState = Pause
		} else if game.currentState == Pause {
			game.currentState = InGame
		} else {
			log.Fatalln("Game state should be InGame or Pause")
		}
	}
}

func getMouseWorldPos() rl.Vector2 {
	worldX := (rl.GetMousePosition().X / float32(rl.GetScreenWidth())) * WorldWidth
	worldY := (rl.GetMousePosition().Y / float32(rl.GetScreenHeight())) * WorldHeight
	worldPos := rl.Vector2{X: worldX, Y: worldY}
	return worldPos
}
