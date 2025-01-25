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

type Game struct {
    ballPosition  rl.Vector2
    ballColor     rl.Color
    targets       []rl.Vector2
    targetRadius  float32
    spawnInterval float32
    lastSpawnTime float32
}

func NewGame() Game {
    return Game{
        ballPosition:  rl.NewVector2(100, 100),
        ballColor:     rl.GetColor(0xf96e61ff),
        targets:       make([]rl.Vector2, 0),
        targetRadius:  30.0,
        spawnInterval: 1.5,
        lastSpawnTime: 0,
    }
}

const ScreenWidth = 800.0
const ScreenHeight = 600.0

func main() {

    rl.InitWindow(ScreenWidth, ScreenHeight, "Click me!")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)
    game := NewGame()

    var timer float32

    for !rl.WindowShouldClose() {
        dt := rl.GetFrameTime()
        timer += dt

        game.ballPosition = rl.GetMousePosition()

        if timer-game.lastSpawnTime >= game.spawnInterval {
            x := rand.Float32() * (ScreenWidth - 2.0*game.targetRadius) + game.targetRadius
            y := rand.Float32() * (ScreenHeight - 2.0*game.targetRadius) + game.targetRadius
            game.targets = append(game.targets, rl.NewVector2(x, y))
            game.lastSpawnTime = timer

            if game.spawnInterval > 0.2 {
                game.spawnInterval -= 0.1 * dt
            }
        }


        if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
            targets := make([]rl.Vector2, len(game.targets))
            index := 0
            for _, target := range game.targets {
                if !rl.CheckCollisionPointCircle(rl.GetMousePosition(),target, game.targetRadius-10){
                    targets[index] = target
                    index++
                }
            }
            game.targets = targets[:index]
        }

        rl.BeginDrawing()

        rl.ClearBackground(rl.GetColor(0xf9d8c2FF))
        rl.DrawCircleV(game.ballPosition, 40, game.ballColor)

        for _, target := range game.targets {
            rl.DrawCircleV(target, float32(game.targetRadius-10), rl.GetColor(0x9f80fcff))
        }
        rl.EndDrawing()
    }
}





