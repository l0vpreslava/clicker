// TODO: Привести код в порядок, нормально расставить кнопочки
// TODO: Оформить вид Game Over
// TODO: Оформить вид Pause
// TODO: Разобраться с поддержкой разных разрешений экрана и fullscreen
// TODO: Пофиксить хитбокс птиц
// TODO: Добавить визуальный эффект выстрела
// TODO: Разобраться в кросс-компиляции попробовать подготовить релиз игры под Макос, Винду и Линукс
// TODO: Выложить игру на itch.io

package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Menu = iota
	InSettings
	InGame
	Pause
	GameOver
)

const (
	Left = -1
	Right = 1
)

const ScreenWidth = 800.0
const ScreenHeight = 600.0
const BannerTimeMax float32 = 0.5
const TargetRadius float32 = 30

type Game struct {
	ballPosition   rl.Vector2
	ballColor      rl.Color
	targets        []Target
	timer          float32
	spawnInterval  float32
	lastSpawnTime  float32
	removedTargets []RemovedTarget
	score          int
	scores         []int
	escapedBird    int
	currentState   int
	settings       Settings
	assets         Assets
	showBanner     bool
	bannerTimer    float32
	shouldClose    bool
	bannerAlpha    float32
}

type Assets struct {
	font          rl.Font
	bgMenu        rl.Texture2D
	bgGame        rl.Texture2D
	banner        rl.Texture2D
	birdWingsUp   rl.Texture2D
	birdWingsDown rl.Texture2D
	birdDead      rl.Texture2D
	sight         rl.Texture2D
	dislike       rl.Texture2D
	selectButton  rl.Sound
	startGame     rl.Sound
	pause         rl.Sound
	bannerSound   rl.Sound
	shoot_1       rl.Sound
	shoot_2       rl.Sound
	birdDied_1    rl.Sound
	birdDied_2    rl.Sound
}

type Settings struct {
	screenResolution  []string
	currentResolution string
	fullscreen        bool
	gameVolume        float32
	scale             float32
}

func NewGame() Game {
	return Game{
		ballPosition:   rl.NewVector2(100, 100),
		ballColor:      rl.GetColor(0xf96e61ff),
		targets:        make([]Target, 0),
		timer:          0,
		spawnInterval:  1.5,
		lastSpawnTime:  0,
		removedTargets: make([]RemovedTarget, 0),
		score:          0,
		scores:         make([]int, 0),
		currentState:   Menu,
		settings:       NewSettigs(),
		assets:         LoadAssets(),
		showBanner:     false,
		bannerTimer:    0,
		shouldClose:    false,
		bannerAlpha:    0,
	}
}

func (game *Game) Destroy() {
	game.assets.Unload()
}

func NewSettigs() Settings {
	return Settings{
		screenResolution:  []string{"800x600", "1200x800", "1440x900"},
		currentResolution: "800x600",
		fullscreen:        false,
		gameVolume:        0.50,
		scale:             1.5,
	}
}

func LoadAssets() Assets {
	return Assets{
		font:          rl.LoadFont("assets/fonts/pixeleum-48.ttf"),
		bgMenu:        rl.LoadTexture("assets/sprites/bg-menu.png"),
		bgGame:        rl.LoadTexture("assets/sprites/bg-game.png"),
		banner:        rl.LoadTexture("assets/sprites/bg-rep.png"),
		dislike:       rl.LoadTexture("assets/sprites/hand.png"),
		birdWingsUp:   rl.LoadTexture("assets/sprites/bird1.png"),
		birdWingsDown: rl.LoadTexture("assets/sprites/bird2.png"),
		birdDead:      rl.LoadTexture("assets/sprites/bird3.png"),
		sight:         rl.LoadTexture("assets/sprites/sight.png"),
		bannerSound:   rl.LoadSound("assets/sounds/banner.wav"),
		birdDied_1:    rl.LoadSound("assets/sounds/bird-died-1.wav"),
		birdDied_2:    rl.LoadSound("assets/sounds/bird-died-2.wav"),
		pause:         rl.LoadSound("assets/sounds/pause.wav"),
		selectButton:  rl.LoadSound("assets/sounds/select-button.wav"),
		shoot_1:       rl.LoadSound("assets/sounds/shoot-1.wav"),
		shoot_2:       rl.LoadSound("assets/sounds/shoot-2.wav"),
		startGame:     rl.LoadSound("assets/sounds/start-game.wav"),
	}
}

func (assets *Assets) Unload() {
	rl.UnloadFont(assets.font)
	rl.UnloadTexture(assets.bgMenu)
	rl.UnloadTexture(assets.bgGame)
	rl.UnloadTexture(assets.banner)
	rl.UnloadTexture(assets.dislike)
	rl.UnloadTexture(assets.birdWingsUp)
	rl.UnloadTexture(assets.birdWingsDown)
	rl.UnloadTexture(assets.birdDead)
	rl.UnloadTexture(assets.sight)
	rl.UnloadSound(assets.selectButton)
	rl.UnloadSound(assets.startGame)
	rl.UnloadSound(assets.pause)
	rl.UnloadSound(assets.bannerSound)
	rl.UnloadSound(assets.shoot_1)
	rl.UnloadSound(assets.shoot_2)
	rl.UnloadSound(assets.birdDied_1)
	rl.UnloadSound(assets.birdDied_2)
}

type Target struct {
	position          rl.Vector2
	velocity          rl.Vector2
	direction         int
	animationFrame    int
	animationTimer    float32
	animationInterval float32
}

type RemovedTarget struct {
	position  rl.Vector2
	timer     float32
	direction int
}

func ShowDeleteAnimation(target RemovedTarget, dt float32) rl.Vector2 {
	const fallSpeed = 500.0    // u
	const acceleration = 200.0 // a

	target.timer += dt // t
	// s = ut + 1/2*at^2
	target.position.Y += fallSpeed*target.timer + 0.5*acceleration*target.timer*target.timer

	return target.position
}

func (game *Game) updateAnimations(dt float32) {
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

func SaveScore(score int, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scoreStr := fmt.Sprintf("%d\n", score)

	_, err = file.WriteString(scoreStr)
	if err != nil {
		return err
	}

	return nil
}

func (game *Game) loadScores(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			game.scores = []int{}
		}
		fmt.Printf("Error loading scores: %s\n", err)
		return
	}

	lines := strings.Split(string(data), "\n")

	scores := make([]int, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		score, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Error parsing scores: %s\n", err)
		}
		scores = append(scores, score)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	game.scores = scores
}

func DrawHighScores(scores []int, font rl.Font, position rl.Vector2, fontSize float32, spacing float32, color rl.Color) {
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

func (game *Game) Update() {
	game.handleInput()
	game.updateState()

	rl.BeginDrawing()

	game.draw()
	game.handleUI()

	rl.EndDrawing()
}

func (game *Game) handleUI() {
	switch game.currentState {
	case Menu:
		rl.DrawTextureEx(
			game.assets.bgMenu,
			rl.Vector2{X: 0, Y: 0},
			0,
			10,
			rl.White)

		DrawHighScores(game.scores, game.assets.font, rl.NewVector2(600, 130), 30, 2, rl.GetColor(0xff0000ff))

		if rg.Button(rl.Rectangle{X: 50, Y: 150, Width: 100, Height: 50}, "START") {
			rl.PlaySound(game.assets.selectButton)
			rl.PlaySound(game.assets.startGame)
			game.currentState = InGame
		}

		if rg.Button(rl.Rectangle{X: 50, Y: 250, Width: 100, Height: 50}, "Settings") {
			rl.PlaySound(game.assets.selectButton)
			game.currentState = InSettings
		}

		if rg.Button(rl.Rectangle{X: 50, Y: 350, Width: 100, Height: 50}, "Quit") {
			rl.PlaySound(game.assets.selectButton)
			game.shouldClose = true
		}

	case InSettings:
		if rl.IsKeyPressed(rl.KeyEscape) || rg.Button(rl.Rectangle{X: 10, Y: 10, Width: 20, Height: 20}, "#114#") {
			rl.PlaySound(game.assets.selectButton)
			game.currentState = Menu
		}

		text := ""
		if game.settings.fullscreen {
			text = "Fullscreen"
		} else {
			text = "Window"
		}

		if rg.Button(rl.Rectangle{X: 300, Y: 50, Width: 200, Height: 50}, text) {
			rl.PlaySound(game.assets.selectButton)
			game.settings.fullscreen = !game.settings.fullscreen
			rl.ToggleFullscreen()
		}

		if rg.Button(rl.Rectangle{X: 300, Y: 150, Width: 200, Height: 50}, game.settings.currentResolution) {
			rl.PlaySound(game.assets.selectButton)
			screenResLen := len(game.settings.screenResolution)

			for i, resolution := range game.settings.screenResolution {
				if resolution == game.settings.currentResolution && (i+1) != screenResLen {
					game.settings.currentResolution = game.settings.screenResolution[i+1]
					SetWindowSize(game.settings.currentResolution)
					break
				} else if i == screenResLen-1 {
					game.settings.currentResolution = game.settings.screenResolution[0]
					SetWindowSize(game.settings.currentResolution)
				}
			}
		}

		game.settings.gameVolume = rg.SliderBar(rl.Rectangle{X: 300, Y: 250, Width: 200, Height: 50}, "Volume", "", game.settings.gameVolume, 0, 1)
		rl.SetMasterVolume(game.settings.gameVolume)

	case InGame:
		text := fmt.Sprintf("Score: %d", game.score)
		rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: 10, Y: 10}, 45, 10, rl.GetColor(0xcb65f7ff))
		rl.DrawTextureV(game.assets.dislike, rl.Vector2{X: 10, Y: 60}, rl.White)
		text = fmt.Sprintf(": %d", game.escapedBird)
		rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: float32(10 + game.assets.dislike.Width), Y: 60}, 45, 10, rl.GetColor(0xf7658eff))

	case Pause:
		text := fmt.Sprintf("Score: %d", game.score)
		rl.DrawTextEx(game.assets.font, text, rl.Vector2{X: 10, Y: 10}, 45, 10, rl.GetColor(0x4d2f1fff))
		text = "Pause"
		pos := CenterText(text, 50, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
		rl.DrawTextEx(game.assets.font, text, pos, 50, 10, rl.GetColor(0x4d2f1fff))

		if rg.Button(rl.Rectangle{X: 50, Y: 350, Width: 100, Height: 50}, "Continue") {
			rl.PlaySound(game.assets.startGame)
			game.currentState = InGame
		}
		if rg.Button(rl.Rectangle{X: 150, Y: 350, Width: 100, Height: 50}, "Quit and save result") {
			rl.PlaySound(game.assets.selectButton)
			err := SaveScore(game.score, "score.txt")
			if err != nil {
				fmt.Println("Error saving score:", err)
			}
			game.shouldClose = true
		}

	case GameOver:
		DrawHighScores(game.scores, game.assets.font, rl.NewVector2(50, 150), 30, 2, rl.GetColor(0x4d2f1fff))
		text := fmt.Sprintf("GAME OVER\n score: %d", game.score)
		pos := CenterText(text, 50, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))

		rl.DrawTextEx(game.assets.font, text, pos, 45, 10, rl.GetColor(0x4d2f1fff))

		if rg.Button(rl.Rectangle{X: 250, Y: 350, Width: 100, Height: 50}, "Retry") {
			rl.PlaySound(game.assets.selectButton)
			game.score = 0
			game.targets = make([]Target, 0)
			game.timer = 0
			game.lastSpawnTime = 0
			game.removedTargets = make([]RemovedTarget, 0)
			game.escapedBird = 0
			rl.PlaySound(game.assets.startGame)
			game.currentState = InGame
		}

		if rg.Button(rl.Rectangle{X: 400, Y: 350, Width: 100, Height: 50}, "Quit") {
			rl.PlaySound(game.assets.selectButton)
			game.shouldClose = true
		}

	}
}

func (game *Game) draw() {
	rl.ClearBackground(rl.GetColor(0x553a7aff))

	switch game.currentState {
	case InGame:
		rl.DrawTextureEx(
			game.assets.bgGame,
			rl.Vector2{X: 0, Y: 0},
			0,
			10,
			rl.White)

		for _, target := range game.targets {
			var birdFrame rl.Texture2D
			if target.animationFrame == 0 {
				birdFrame = game.assets.birdWingsUp
			} else {
				birdFrame = game.assets.birdWingsDown
			}

			if target.direction == Left {
				flipRec := rl.NewRectangle(
					float32(birdFrame.Width),
					0,
					-float32(birdFrame.Width),
					float32(birdFrame.Height))
				rl.DrawTextureRec(
					birdFrame,
					flipRec,
					rl.Vector2{X: target.position.X - TargetRadius, Y: target.position.Y - TargetRadius - 10},
					rl.White)
			} else if target.direction == Right {
				rl.DrawTextureV(
					birdFrame,
					rl.Vector2{X: target.position.X - TargetRadius, Y: target.position.Y - TargetRadius - 10},
					rl.White)
			} else {
				log.Fatalln("target.direction should be Left or Right")
			}
		}

		for _, removedTarget := range game.removedTargets {
			if removedTarget.direction == -1 {
				flipRec := rl.NewRectangle(
					float32(game.assets.birdDead.Width),
					0,
					-float32(game.assets.birdDead.Width),
					float32(game.assets.birdDead.Height))
				rl.DrawTextureRec(
					game.assets.birdDead,
					flipRec,
					rl.Vector2{X: removedTarget.position.X - TargetRadius, Y: removedTarget.position.Y - TargetRadius - 10},
					rl.White)
			} else {
				rl.DrawTextureV(
					game.assets.birdDead,
					rl.Vector2{X: removedTarget.position.X - TargetRadius, Y: removedTarget.position.Y - TargetRadius - 10},
					rl.White)
			}
		}

		if game.showBanner {
			color := rl.ColorAlpha(rl.White, game.bannerAlpha)
			rl.DrawTextureEx(game.assets.banner, rl.Vector2{X: 100, Y: 100}, 0, 10, color)
		}

		rl.DrawTextureEx(
			game.assets.sight,
			rl.Vector2{X: game.ballPosition.X - 32*game.settings.scale - 1, Y: game.ballPosition.Y - 32*game.settings.scale},
			0,
			game.settings.scale,
			rl.White)
	}
}

func (game *Game) updateState() {
	dt := rl.GetFrameTime()
	game.timer += dt

	if game.currentState == InGame {
		rl.HideCursor()
	} else {
		rl.ShowCursor()
	}

	switch game.currentState {
	case InGame:
		game.spawnBirds(dt)
		game.updateBirds(dt)
		game.updateBanner(dt)
		game.updateAnimations(dt)
	}
}

func (game *Game) saveScore() {
	file, err := os.OpenFile("score.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	for _, score := range game.scores {
		scoreStr := fmt.Sprintf("%d\n", score)
		_, err = file.WriteString(scoreStr)
		if err != nil {
			fmt.Printf("Error writing file: %s\n", err)
			return
		}
	}
}

func (game *Game) updateBanner(dt float32) {
	if game.showBanner {
		game.bannerTimer -= dt
		if game.bannerTimer <= 0 {
			game.showBanner = false
		}
	}

	game.bannerAlpha = game.bannerTimer / BannerTimeMax
}

func (game *Game) updateBirds(dt float32) {
	for i := range game.targets {
		game.targets[i].position = rl.Vector2Add(game.targets[i].position, rl.Vector2Scale(game.targets[i].velocity, dt))

		game.targets[i].animationTimer += dt
		if game.targets[i].animationTimer >= game.targets[i].animationInterval {
			game.targets[i].animationFrame = (game.targets[i].animationFrame + 1) % 2
			game.targets[i].animationTimer = 0
		}

	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		var sound rl.Sound
		if rand.IntN(2) == 0 {
			sound = game.assets.shoot_1
		} else {
			sound = game.assets.shoot_2
		}
		rl.PlaySound(sound)

		targets := make([]Target, len(game.targets))
		index := 0
		for _, target := range game.targets {
			if !rl.CheckCollisionPointCircle(rl.GetMousePosition(), target.position, TargetRadius-10) {
				targets[index] = target
				index++
			} else {
				game.removedTargets = append(game.removedTargets, RemovedTarget{position: target.position, timer: 0, direction: target.direction})
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

		game.targets = targets[:index]
	}

	var newTargets []Target
	for _, target := range game.targets {
		if (target.direction == Right && target.position.X > float32(rl.GetScreenWidth())+TargetRadius) ||
			(target.direction == Left && target.position.X < float32(0-TargetRadius)) {
			game.escapedBird++
			rl.PlaySound(game.assets.bannerSound)
			game.showBanner = !game.showBanner
			game.bannerTimer = BannerTimeMax
		} else {
			newTargets = append(newTargets, target)
		}
	}

	game.targets = newTargets

	if game.escapedBird >= 5 {
		for _, target := range game.targets {
			game.removedTargets = append(game.removedTargets, RemovedTarget{position: target.position, timer: 0})
		}
		game.targets = make([]Target, 0)
		game.scores = append(game.scores, game.score)
		sort.Sort(sort.Reverse(sort.IntSlice(game.scores)))
		go game.saveScore()
		game.currentState = GameOver
	}
}

func (game *Game) spawnBirds(dt float32) {
	if game.timer-game.lastSpawnTime >= game.spawnInterval {
		var x, y float32
		var velocity rl.Vector2
		var direction int
		y = rand.Float32()*(ScreenHeight-2.0*TargetRadius) + TargetRadius

		switch rand.IntN(2) {
		case 0:
			x = -TargetRadius
			velocity = rl.NewVector2(1, 0)
			direction = Right

		case 1:
			x = ScreenWidth + TargetRadius
			velocity = rl.NewVector2(-1, 0)
			direction = Left
		}

		speed := rand.Float32()*300 + 100
		velocity = rl.Vector2Scale(velocity, speed)

		game.targets = append(game.targets, Target{
			position:          rl.NewVector2(x, y),
			velocity:          velocity,
			direction:         direction,
			animationFrame:    0,
			animationTimer:    0,
			animationInterval: 0.1,
		})

		game.lastSpawnTime = game.timer

		if game.spawnInterval > 0.5 {
			game.spawnInterval -= 0.1 * dt
		}
	}
}

func (game *Game) handleInput() {
	if rl.WindowShouldClose() {
		game.shouldClose = true
	}

	game.ballPosition = rl.GetMousePosition()
	
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

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyF5)

	game := NewGame()
	defer game.Destroy()

	game.loadScores("score.txt")

	for !game.shouldClose {
		game.Update()
	}
}
