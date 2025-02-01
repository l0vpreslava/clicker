//TODO: Нарисовать фоны, мишени, курсор, эффект курсора при нажатии
//TODO: Интегрировать нарисованные ресурсы в игру
//TODO: Начать сохранять рекорды по очкам в файл
//TODO: Добавить в меню таблицу рекордов
//TODO: Добавить в игру звуковое сопровождение: эффекты нажатия кнопок в меню, эффект нажатия курсора в игре, эффект уничтожения мишени
//TODO: Разобраться в кросс-компиляции попробовать подготовить релиз игры под Макос, Винду и Линукс
//TODO: Выложить игру на itch.io

package main

import (
    "fmt"
    rg "github.com/gen2brain/raylib-go/raygui"
    rl "github.com/gen2brain/raylib-go/raylib"
    "math/rand/v2"
    "strconv"
    "strings"
)

const (
    Menu = iota
    InSettings
    InGame
    Pause
    GameOver
)

const ScreenWidth = 800.0
const ScreenHeight = 600.0

type Game struct {
    ballPosition   rl.Vector2
    ballColor      rl.Color
    targets        []Target
    targetRadius   float32
    timer          float32
    spawnInterval  float32
    lastSpawnTime  float32
    removedTargets []RemovedTarget
    score          int
    font           rl.Font
    currentState   int
}

type Settings struct {
    screenResolution  []string
    currentResolution string
    fullscreen        bool
    gameVolume        float32
}

func NewGame() Game {
    return Game{
        ballPosition:   rl.NewVector2(100, 100),
        ballColor:      rl.GetColor(0xf96e61ff),
        targets:        make([]Target, 0),
        targetRadius:   30.0,
        timer:          0,
        spawnInterval:  1.5,
        lastSpawnTime:  0,
        removedTargets: make([]RemovedTarget, 0),
        score:          0,
        font:           rl.LoadFont("assets/fonts/pixeleum-48.ttf"),
        currentState:   Menu,
    }
}

func NewSettigs() Settings {
    return Settings{
        screenResolution:  []string{"800x600", "1200x800", "1440x900"},
        currentResolution: "800x600",
        fullscreen:        false,
        gameVolume:        0.50,
    }
}

type Target struct {
    position rl.Vector2
    velocity rl.Vector2
}

type RemovedTarget struct {
    position rl.Vector2
    timer    float32
}

func ShowDeleteAnimation(target RemovedTarget, dt float32) rl.Vector2 {
    const fallSpeed = 500.0    // u
    const acceleration = 200.0 // a

    target.timer += dt // t
    // s = ut + 1/2*at^2
    target.position.Y += fallSpeed*target.timer + 0.5*acceleration*target.timer*target.timer

    return target.position
}

func UpdateAnimation(game *Game, dt float32) {
    for i := range game.removedTargets {
        game.removedTargets[i].position = ShowDeleteAnimation(game.removedTargets[i], dt)
    }

    newRemovedTargets := make([]RemovedTarget, 0)
    for _, removedTarget := range game.removedTargets {
        if removedTarget.position.Y < float32(rl.GetScreenHeight()) {
            newRemovedTargets = append(newRemovedTargets, removedTarget)
        }
    }
    game.removedTargets = newRemovedTargets
}

func SetWindowSize(resolution string) {
    parts := strings.Split(resolution, "x")
    if len(parts) != 2 {
        fmt.Println("Invalid resolution format.")
        return
    }

    width, err := strconv.Atoi(parts[0])
    if err != nil {
        fmt.Println("Error converting width:", err)
        return
    }

    height, err := strconv.Atoi(parts[1])
    if err != nil {
        fmt.Println("Error converting height:", err)
        return
    }

    rl.SetWindowSize(width, height)
}

func Update(game *Game, settings *Settings) {
    dt := rl.GetFrameTime()
    game.timer += dt

    game.ballPosition = rl.GetMousePosition()

    switch game.currentState {
    case Menu:
        if rg.Button(rl.Rectangle{X: 50, Y: 150, Width: 100, Height: 50}, "START") {
            game.currentState = InGame
        }

        if rg.Button(rl.Rectangle{X: 50, Y: 250, Width: 100, Height: 50}, "Settings") {
            game.currentState = InSettings
        }

        if rg.Button(rl.Rectangle{X: 50, Y: 350, Width: 100, Height: 50}, "Quit") {
            rl.CloseWindow()
        }
    case InSettings:
        if rl.IsKeyPressed(rl.KeyEscape) || rg.Button(rl.Rectangle{X: 10, Y: 10, Width: 20, Height: 20}, "#114#") {
            game.currentState = Menu
        }

        text := ""
        if settings.fullscreen {
            text = "Fullscreen"
        } else {
            text = "Window"
        }

        if rg.Button(rl.Rectangle{X: 300, Y: 50, Width: 200, Height: 50}, text) {
            settings.fullscreen = !settings.fullscreen
            rl.ToggleFullscreen()
        }

        if rg.Button(rl.Rectangle{X: 300, Y: 150, Width: 200, Height: 50}, settings.currentResolution) {
            screenResLen := len(settings.screenResolution)

            for i, resolution := range settings.screenResolution {
                if resolution == settings.currentResolution && (i+1) != screenResLen {
                    settings.currentResolution = settings.screenResolution[i+1]
                    SetWindowSize(settings.currentResolution)
                    break
                } else if i == screenResLen-1 {
                    settings.currentResolution = settings.screenResolution[0]
                    SetWindowSize(settings.currentResolution)
                }
            }
        }

        settings.gameVolume = rg.SliderBar(rl.Rectangle{X: 300, Y: 250, Width: 200, Height: 50}, "Volume", "", settings.gameVolume, 0, 1)

    case InGame:
        if rl.IsKeyDown(rl.KeyEscape) {
            game.currentState = Pause
        }

        if game.timer-game.lastSpawnTime >= game.spawnInterval {
            side := rand.IntN(2)

            var x, y float32
            var velocity rl.Vector2
            y = rand.Float32()*(ScreenHeight-2.0*game.targetRadius) + game.targetRadius

            switch side {
            case 0: // Left
                x = -game.targetRadius
                velocity = rl.NewVector2(1, 0) // Move right
            case 1: // Right
                x = ScreenWidth + game.targetRadius
                velocity = rl.NewVector2(-1, 0) // Move left
            }

            speed := rand.Float32()*300 + 100
            velocity = rl.Vector2Scale(velocity, speed)

            game.targets = append(game.targets, Target{
                position: rl.NewVector2(x, y),
                velocity: velocity,
            })

            game.lastSpawnTime = game.timer

            if game.spawnInterval > 0.5 {
                game.spawnInterval -= 0.1 * dt
            }
        }

        for i := range game.targets {
            game.targets[i].position = rl.Vector2Add(game.targets[i].position, rl.Vector2Scale(game.targets[i].velocity, dt))
        }

        if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
            targets := make([]Target, len(game.targets))
            index := 0
            for _, target := range game.targets {
                if !rl.CheckCollisionPointCircle(rl.GetMousePosition(), target.position, game.targetRadius-10) {
                    targets[index] = target
                    index++
                } else {
                    game.removedTargets = append(game.removedTargets, RemovedTarget{position: target.position, timer: 0})
                    game.score++
                }
            }
            game.targets = targets[:index]
        }

        UpdateAnimation(game, dt)

        if len(game.targets) >= 15 {
            for _, target := range game.targets {
                game.removedTargets = append(game.removedTargets, RemovedTarget{position: target.position, timer: 0})
            }
            game.targets = make([]Target, 0)
            game.currentState = GameOver
        }
    case Pause:
        if rg.Button(rl.Rectangle{X: 50, Y: 350, Width: 100, Height: 50}, "Continue") {
            game.currentState = InGame
        }
        if rg.Button(rl.Rectangle{X: 150, Y: 350, Width: 100, Height: 50}, "Quit and save result") {
            rl.CloseWindow()
        }
    case GameOver:
        UpdateAnimation(game, dt)

        if rg.Button(rl.Rectangle{X: 250, Y: 350, Width: 100, Height: 50}, "Retry") {
            rl.CloseWindow()
        }

        if rg.Button(rl.Rectangle{X: 400, Y: 350, Width: 100, Height: 50}, "Quit") {
            rl.CloseWindow()
        }
    }

    rl.BeginDrawing()

    rl.ClearBackground(rl.GetColor(0xf9d8c2FF))
    rl.DrawCircleV(game.ballPosition, 40, game.ballColor)
    for _, target := range game.targets {
        rl.DrawCircleV(target.position, float32(game.targetRadius-10), rl.GetColor(0x9f80fcff))
    }

    for _, removedTarget := range game.removedTargets {
        rl.DrawCircleV(removedTarget.position, float32(game.targetRadius-10), rl.GetColor(0x9f80fcff))
    }

    switch game.currentState {
    case Menu:

    case InGame:

        text := fmt.Sprintf("Score: %d", game.score)
        rl.DrawTextEx(game.font, text, rl.Vector2{X: 10, Y: 10}, 45, 10, rl.GetColor(0x4d2f1fff))

    case Pause:
        text := fmt.Sprintf("Score: %d", game.score)
        rl.DrawTextEx(game.font, text, rl.Vector2{X: 10, Y: 10}, 45, 10, rl.GetColor(0x4d2f1fff))
        text = "Pause"
        pos := CenterText(text, 50, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
        rl.DrawTextEx(game.font, text, pos, 50, 10, rl.GetColor(0x4d2f1fff))

    case GameOver:
        text := fmt.Sprintf("GAME OVER\n score: %d", game.score)
        pos := CenterText(text, 50, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))

        rl.DrawTextEx(game.font, text, pos, 45, 10, rl.GetColor(0x4d2f1fff))

    }
    rl.EndDrawing()
}

func CenterText(text string, fontSize int32, screenWidth int32, screenHeight int32) rl.Vector2 {
    textWidth := rl.MeasureText(text, fontSize)
    textHeight := fontSize

    posX := float32((screenWidth - textWidth) / 2)
    posY := float32((screenHeight - textHeight) / 2)

    return rl.Vector2{X: posX, Y: posY}
}

func main() {

    rl.InitWindow(ScreenWidth, ScreenHeight, "Click me!")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)
    rl.SetExitKey(rl.KeyF5)
    rl.SetWindowState(rl.FlagWindowResizable)

    game := NewGame()
    settings := NewSettigs()

    for !rl.WindowShouldClose() {

        Update(&game, &settings)
    }
}
