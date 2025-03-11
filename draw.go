package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (game *Game) draw() {
	rl.BeginTextureMode(game.RenderTexture)
	rl.ClearBackground(rl.GetColor(0x553a7aff))

	switch game.currentState {
	case InGame:
		game.drawBackground()
		game.drawBirds()

		if game.showBanner {
			game.drawBanner()
		}

		game.drawCursor()

	case Pause:
		game.drawBackground()
	}

	rl.EndTextureMode()

	game.drawRenderTextureToTheScreen()
}

func (game *Game) drawCursor() {
	rl.DrawTextureEx(
		game.assets.sight,
		rl.Vector2SubtractValue(getMouseWorldPos(), 32),
		0,
		1,
		rl.White)
}

func (game *Game) drawRenderTextureToTheScreen() {
	sourceRec := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(game.RenderTexture.Texture.Width),
		Height: -float32(game.RenderTexture.Texture.Height),
	}
	scale := PossibleScales[game.settings.ScaleIndex]
	destRec := rl.Rectangle{
		X:      -float32(scale),
		Y:      -float32(scale),
		Width:  float32(rl.GetScreenWidth()) + scale*2,
		Height: float32(rl.GetScreenHeight()) + scale*2,
	}
	rl.DrawTexturePro(
		game.RenderTexture.Texture,
		sourceRec,
		destRec,
		rl.Vector2{},
		0,
		rl.White)
}

func (game *Game) drawBirds() {
	game.drawAliveBirds()
	game.drawKilledBirds()
}

func (game *Game) drawAliveBirds() {
	for _, bird := range game.birds {
		var birdFrame rl.Texture2D
		if bird.animationFrame == 0 {
			birdFrame = game.assets.birdWingsUp
		} else {
			birdFrame = game.assets.birdWingsDown
		}

		if bird.direction == Left {
			flipRec := rl.NewRectangle(
				float32(birdFrame.Width),
				0,
				-float32(birdFrame.Width),
				float32(birdFrame.Height))
			rl.DrawTextureRec(
				birdFrame,
				flipRec,
				rl.Vector2{X: bird.position.X - BirdRadius, Y: bird.position.Y - BirdRadius - 10},
				rl.White)
		} else if bird.direction == Right {
			rl.DrawTextureV(
				birdFrame,
				rl.Vector2{X: bird.position.X - BirdRadius, Y: bird.position.Y - BirdRadius - 10},
				rl.White)
		} else {
			log.Fatalln("bird.direction should be Left or Right")
		}
	}
}

func (game *Game) drawKilledBirds() {
	for _, bird := range game.killedBirds {
		if bird.direction == -1 {
			flipRec := rl.NewRectangle(
				float32(game.assets.birdDead.Width),
				0,
				-float32(game.assets.birdDead.Width),
				float32(game.assets.birdDead.Height))
			rl.DrawTextureRec(
				game.assets.birdDead,
				flipRec,
				rl.Vector2{X: bird.position.X - BirdRadius, Y: bird.position.Y - BirdRadius - 10},
				rl.White)
		} else {
			rl.DrawTextureV(
				game.assets.birdDead,
				rl.Vector2{X: bird.position.X - BirdRadius, Y: bird.position.Y - BirdRadius - 10},
				rl.White)
		}
	}
}

func (game *Game) drawBanner() {
	color := rl.ColorAlpha(rl.White, game.bannerAlpha)
	var padding float32 = 50.0
	renderTextureSize := getTextureSize(game.RenderTexture.Texture)
	destRec := rl.Rectangle{
		X:      float32(padding),
		Y:      float32(padding),
		Width:  renderTextureSize.Width - 2.0*padding,
		Height: renderTextureSize.Height - 2.0*padding,
	}
	rl.DrawTexturePro(
		game.assets.banner,
		getTextureSize(game.assets.banner),
		destRec,
		rl.Vector2{},
		0,
		color)
}

func (game *Game) drawBackground() {
	rl.DrawTexturePro(
		game.assets.bgGame,
		getTextureSize(game.assets.bgGame),
		getTextureSize(game.RenderTexture.Texture),
		rl.Vector2{},
		0,
		rl.White)
}
