package main

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (game *Game) handleUI() {
	centreScreenWidth := rl.GetScreenWidth() / 2
	centreScreenHeight := rl.GetScreenHeight() / 2
	xButton := float32(centreScreenWidth) - 250/2

	switch game.currentState {
	case Menu:
		game.handleMenu(centreScreenWidth, centreScreenHeight)
	case InSettings:
		game.handleSettings(xButton, centreScreenHeight)
	case InGame:
		game.handleInGame()
	case Pause:
		game.handlePause(xButton, centreScreenHeight)
	case GameOver:
		game.handleGameOver(xButton, centreScreenHeight)
	}
}

func (game *Game) handleMenu(centerScreenWidth int, centerScreenHeight int) {
	rl.DrawTexturePro(
		game.assets.bgMenu,
		getTextureSize(game.assets.bgMenu),
		rl.Rectangle{Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())},
		rl.Vector2{},
		0,
		rl.White)

	drawHighScores(
		game.scores,
		game.assets.font,
		rl.NewVector2(float32(rl.GetScreenWidth())-300, float32(160)),
		32,
		2,
		game.TextColor)

	if rg.Button(rl.Rectangle{X: float32(centerScreenWidth) / 2, Y: float32(centerScreenHeight) - 100, Width: 200, Height: 50}, "START") {
		rl.PlaySound(game.assets.selectButton)
		rl.PlaySound(game.assets.startGame)
		game.currentState = InGame
	}

	if rg.Button(rl.Rectangle{X: float32(centerScreenWidth) / 2, Y: float32(centerScreenHeight), Width: 200, Height: 50}, "Settings") {
		rl.PlaySound(game.assets.selectButton)
		game.currentState = InSettings
	}

	if rg.Button(rl.Rectangle{X: float32(centerScreenWidth) / 2, Y: float32(centerScreenHeight) + 100, Width: 200, Height: 50}, "Quit") {
		rl.PlaySound(game.assets.selectButton)
		game.shouldClose = true
	}
}

func (game *Game) handleSettings(xButton float32, centerScreenHeight int) {
	if rl.IsKeyPressed(rl.KeyEscape) || rg.Button(rl.Rectangle{X: 10, Y: 10, Width: 20, Height: 20}, "#114#") {
		rl.PlaySound(game.assets.selectButton)
		go game.settings.Save()
		game.currentState = Menu
	}
	if rg.Button(rl.Rectangle{X: xButton, Y: float32(centerScreenHeight) - 200, Width: 250, Height: 50}, "Toggle Borderless") {
		rl.PlaySound(game.assets.selectButton)
		game.settings.Borderless = !game.settings.Borderless
		rl.ToggleBorderlessWindowed()
	}

	currentResolution := fmt.Sprintf("%dx%d", rl.GetScreenWidth(), rl.GetScreenHeight())
	if rg.Button(rl.Rectangle{X: xButton, Y: float32(centerScreenHeight) - 100, Width: 250, Height: 50}, currentResolution) {
		rl.PlaySound(game.assets.selectButton)
		i := &game.settings.ScaleIndex
		*i = (*i + 1) % len(PossibleScales)
		w, h := GetResolution(*i)
		rl.SetWindowSize(int(w), int(h))

	}

	game.settings.GameVolume = rg.SliderBar(rl.Rectangle{X: xButton, Y: float32(centerScreenHeight), Width: 250, Height: 50}, "Volume", "", game.settings.GameVolume, 0, 1)
	rl.SetMasterVolume(game.settings.GameVolume)
}

func (game *Game) handleInGame() {
	text := fmt.Sprintf("Score: %d", game.score)
	rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: 10, Y: 10}, 48, 10, game.TextColor)
	rl.DrawTextureV(game.assets.dislike, rl.Vector2{X: 10, Y: 60}, rl.White)
	text = fmt.Sprintf(": %d", game.escapedBird)
	rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: float32(10 + game.assets.dislike.Width), Y: 60}, 48, 10, game.TextColor)
}

func (game *Game) handlePause(xButton float32, centreScreenHeight int) {
	text := fmt.Sprintf("Score: %d", game.score)
	rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: 10, Y: 10}, 48, 10, game.TextColor)
	text = "Pause"
	pos := CenterText(text, 50, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
	rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: pos.X, Y: pos.Y - 150}, 48, 10, game.TextColor)

	if rg.Button(rl.Rectangle{X: xButton, Y: float32(centreScreenHeight) - 100, Width: 200, Height: 50}, "Continue") {
		rl.PlaySound(game.assets.startGame)
		game.currentState = InGame
	}

	if rg.Button(rl.Rectangle{X: xButton, Y: float32(centreScreenHeight), Width: 200, Height: 50}, "#185# Menu") {
		rl.PlaySound(game.assets.selectButton)
		game.appendCurrentScore()
		game.score = 0
		game.currentState = Menu
	}
	if rg.Button(rl.Rectangle{X: xButton, Y: float32(centreScreenHeight) + 100, Width: 200, Height: 50}, "#113# Quit") {
		rl.PlaySound(game.assets.selectButton)
		game.appendCurrentScore()
		game.shouldClose = true
	}
}

func (game *Game) handleGameOver(xButton float32, centreScreenHeight int) {
	drawHighScores(game.scores, game.assets.font, rl.NewVector2(50, 150), 32, 2, game.TextColor)

	text := fmt.Sprintf("GAME OVER\n score: %d", game.score)
	pos := CenterText(text, 48, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))

	rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: pos.X, Y: pos.Y - 150}, 48, 10, game.TextColor)

	if rg.Button(rl.Rectangle{X: xButton - 100, Y: float32(centreScreenHeight), Width: 200, Height: 50}, "Retry") {
		rl.PlaySound(game.assets.selectButton)
		game.score = 0
		game.birds = make([]Bird, 0)
		game.timer = 0
		game.lastSpawnTime = 0
		game.killedBirds = make([]KilledBird, 0)
		game.escapedBird = 0
		rl.PlaySound(game.assets.startGame)
		game.appendCurrentScore()
		game.score = 0
		game.currentState = InGame
	}

	if rg.Button(rl.Rectangle{X: xButton + 150, Y: float32(centreScreenHeight), Width: 200, Height: 50}, "Quit") {
		rl.PlaySound(game.assets.selectButton)
		game.appendCurrentScore()
		game.shouldClose = true
	}
}

func drawHighScores(scores []int, font rl.Font, position rl.Vector2, fontSize float32, spacing float32, color rl.Color) {
	y := position.Y

	title := "High Scores:"
	rl.DrawTextEx(font, title, rl.Vector2{X: position.X, Y: y}, fontSize, spacing, color)
	y += fontSize + 10
	for i, score := range scores {
		text := fmt.Sprintf("%d. %d", i+1, score)

		rl.DrawTextEx(font, text, rl.Vector2{X: position.X, Y: y}, fontSize, spacing, color)

		y += fontSize + 5
		if i == 9 {
			return
		}
	}
}

func CenterText(text string, fontSize int32, screenWidth int32, screenHeight int32) rl.Vector2 {
	textWidth := rl.MeasureText(text, fontSize)
	textHeight := fontSize

	posX := float32((screenWidth - textWidth) / 2)
	posY := float32((screenHeight - textHeight) / 2)

	return rl.Vector2{X: posX, Y: posY}
}
