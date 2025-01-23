//TODO: Обработать нажатие мышки так, что если игрок нажал в пределах мишени, то эту мишень удалить из списка
//TODO: Добавить какой-нибудь эффект при уничтожении мишени
//TODO: Добавить счётчик очков
//TODO: Добавить ускорение спавна мишеней со временем
//TODO: Добавить условие проигрыша, когда мишеней станет слишком много
//TODO: Добавить меню, в котором будет кнопка играть, настройки и выход
//TODO: Нарисовать фоны, мишени, курсор, эффект курсора при нажатии
//TODO: Интегрировать нарисованные ресурсы в игру
//TODO: Начать сохранять рекорды по очкам в файл
//TODO: Добавить в меню таблицу рекордов
//TODO: Добавить в игру звуковое сопровождение: эффекты нажатия кнопок в меню, эффект нажатия курсора в игре, эффект уничтожения мишени
//TODO: Разобраться в кросс-компиляции попробовать подготовить релиз игры под Макос, Винду и Линукс
//TODO: Выложить игру на itch.io

package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
    "math/rand/v2"
)

func main() {
    const screenWidth = 800
    const screenHeight = 600

    rl.InitWindow(screenWidth, screenHeight, "Click me!")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)
    ballPosition := rl.NewVector2(100, 100)
    ballColor := rl.GetColor(0xf96e61ff)

    targets := []rl.Vector2{}
    targetRadius := 30

    var timer float32
    spawnInterval := float32(1.5)
    lastSpawnTime := float32(0)

    for !rl.WindowShouldClose() {
        timer += rl.GetFrameTime()

        if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
            ballPosition = rl.GetMousePosition()
        }

        if timer-lastSpawnTime >= spawnInterval {
            x := float32((rand.IntN(screenWidth-2*targetRadius)) + targetRadius)
            y := float32((rand.IntN(screenHeight-2*targetRadius)) + targetRadius)
            targets = append(targets, rl.NewVector2(x, y))
            lastSpawnTime = timer

            if spawnInterval > 0.2 {
                spawnInterval -= 0.0001
            }
        }
        rl.BeginDrawing()

        rl.ClearBackground(rl.GetColor(0xf9d8c2FF))
        rl.DrawCircleV(ballPosition, 40, ballColor)

        for _, target := range targets {
            rl.DrawCircleV(target, float32(targetRadius-10), rl.GetColor(0x9f80fcff))
        }
        rl.EndDrawing()
    }
}
