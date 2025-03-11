package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func getTextureSize(texture rl.Texture2D) rl.Rectangle {
	return rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}
}

func GetGameCacheDir() string {
	cacheDir, err := os.UserConfigDir()
	if err == nil {
		return fmt.Sprintf("%s/flyordie", cacheDir)
	} else {
		return "."
	}
}
