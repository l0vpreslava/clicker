package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
    "math/rand/v2"
)

func main() {
    const screenWidth = 800
    const screenHeight = 600


    rl.InitWindow(screenWidth,screenHeight,"Click me!")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)
    ballPosition := rl.NewVector2(100,100)
    ballColor := rl.GetColor(0xf96e61ff)

    targets := []rl.Vector2{}
    targetRadius := 30


    var timer float32
    spawnInterval := 1.5 
    lastSpawnTime := float32(0)

    for !rl.WindowShouldClose(){
        timer += rl.GetFrameTime()

        if rl.IsMouseButtonPressed(rl.MouseButtonLeft){
            ballPosition = rl.GetMousePosition()
        }

        if timer-lastSpawnTime >= float32(spawnInterval) {
            x := float32((rand.IntN(screenWidth-2*targetRadius)) + targetRadius)
            y := float32((rand.IntN(screenHeight-2*targetRadius)) + targetRadius)
            targets = append(targets, rl.NewVector2(x,y))
            lastSpawnTime = timer

            if spawnInterval > 0.2 {
                spawnInterval -= 0.0001
            }
        }

        rl.BeginDrawing()

        rl.ClearBackground(rl.GetColor(0xf9d8c2FF))
        rl.DrawCircleV(ballPosition, 40, ballColor)
        for _, target := range targets{
            rl.DrawCircleV(target, float32(targetRadius-10), rl.GetColor(0x9f80fcff))
        }
        
        rl.EndDrawing()
    }
}

