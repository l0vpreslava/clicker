// TODO: Разобраться в кросс-компиляции попробовать подготовить релиз игры под Макос, Винду и Линукс
// TODO: Выложить игру на itch.io

package main

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	settings := LoadSettings()

	w, h := GetResolution(settings.ScaleIndex)
	rl.InitWindow(w, h, "Click me!")
	defer rl.CloseWindow()

	if settings.Borderless {
		rl.ToggleBorderlessWindowed()
	}

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetMasterVolume(settings.GameVolume)

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyF5)
	rg.SetStyle(0, rg.TEXT_SIZE, 25)

	game := NewGame()
	defer game.Destroy()

	game.loadScores()
	game.settings = settings

	for !game.shouldClose {
		game.Update()
	}
}
