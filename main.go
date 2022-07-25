package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	scale            = 18
	gridWidth        = 30
	gridHeight       = 30
	screenWidth      = gridWidth * scale
	screenHeight     = gridHeight * scale
	startingTickRate = 0.25
	appleTickLimit   = (gridWidth + gridHeight) * 0.8
	speedLimit       = 0.060
	speedRate        = 0.005
	scoreMultiplier  = 10
)

type Direction int

type Apple struct {
	pos    *rl.Vector2
	ticker int
}

type Snake struct {
	body          []*rl.Vector2
	timer         float32
	tickRate      float32
	isDead        bool
	direction     Direction
	nextDirection Direction
}

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

var (
	snake *Snake
	rect  = rl.Rectangle{
		Width:  scale,
		Height: scale,
		X:      0,
		Y:      0,
	}
	rotation float32
	origin   = rl.Vector2{
		X: 0,
		Y: 0,
	}
	vSync   = false
	showFPS = false
	apple   *Apple
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go Snake")
	setup()
	for !rl.WindowShouldClose() {
		update()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		draw()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func setup() {
	snake = &Snake{
		body:          make([]*rl.Vector2, 0),
		timer:         0.0,
		tickRate:      startingTickRate,
		isDead:        false,
		direction:     Right,
		nextDirection: Right,
	}
	snake.body = append(snake.body, &rl.Vector2{
		X: (gridWidth / 2) * scale,
		Y: (gridHeight / 2) * scale,
	})
	snake.body = append(snake.body, &rl.Vector2{
		X: ((gridWidth / 2) - 1) * scale,
		Y: (gridHeight / 2) * scale,
	})
	apple = &Apple{
		pos: &rl.Vector2{
			X: float32(rl.GetRandomValue(1, gridWidth-1) * scale),
			Y: float32(rl.GetRandomValue(1, gridHeight-1) * scale),
		},
		ticker: 0,
	}
	rotation = 0.0
}

func update() {

	if rl.IsKeyDown(rl.KeyD) && snake.direction != Left {
		snake.nextDirection = Right
	}
	if rl.IsKeyDown(rl.KeyA) && snake.direction != Right {
		snake.nextDirection = Left
	}
	if rl.IsKeyDown(rl.KeyW) && snake.direction != Down {
		snake.nextDirection = Up
	}
	if rl.IsKeyDown(rl.KeyS) && snake.direction != Up {
		snake.nextDirection = Down
	}
	if rl.IsKeyPressed(rl.KeyV) {
		if vSync {
			vSync = !vSync
			rl.SetTargetFPS(0)
		} else {
			vSync = !vSync
			rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
		}
	}
	if rl.IsKeyPressed(rl.KeyF) {
		showFPS = !showFPS
	}
	if rl.IsKeyPressed(rl.KeyR) {
		setup()
		return
	}

	dt := rl.GetFrameTime()
	if snake.isDead {
		rotation += 60.0 * dt
		if rotation > 360 {
			rotation = 0
		}
		return
	}

	snake.timer += dt
	if snake.timer < snake.tickRate {
		return
	}
	snake.timer = 0.0
	snake.direction = snake.nextDirection

	if snake.body[0].X == apple.pos.X && snake.body[0].Y == apple.pos.Y {
		snake.body = append(snake.body, &rl.Vector2{
			X: snake.body[len(snake.body)-1].X,
			Y: snake.body[len(snake.body)-1].Y,
		})
		apple.pos.X = float32(rl.GetRandomValue(1, gridWidth-1) * scale)
		apple.pos.Y = float32(rl.GetRandomValue(1, gridHeight-1) * scale)
		apple.ticker = 0
		if snake.tickRate > speedLimit {
			snake.tickRate -= speedRate
		}
	}

	apple.ticker += 1
	if apple.ticker > appleTickLimit {
		apple.pos.X = float32(rl.GetRandomValue(1, gridWidth-1) * scale)
		apple.pos.Y = float32(rl.GetRandomValue(1, gridHeight-1) * scale)
		apple.ticker = 0
	}

	prev := rl.Vector2{
		X: snake.body[0].X,
		Y: snake.body[0].Y,
	}
	hold := rl.Vector2{
		X: snake.body[0].X,
		Y: snake.body[0].Y,
	}
	for i, vec := range snake.body {
		if i == 0 {
			switch snake.direction {
			case Up:
				vec.Y -= scale
			case Down:
				vec.Y += scale
			case Left:
				vec.X -= scale
			case Right:
				vec.X += scale
			}

			if snake.body[0].X > (gridWidth-1)*scale {
				snake.body[0].X = 0
			}
			if snake.body[0].X < 0 {
				snake.body[0].X = (gridWidth - 1) * scale
			}
			if snake.body[0].Y > (gridHeight-1)*scale {
				snake.body[0].Y = scale
			}
			if snake.body[0].Y < scale {
				snake.body[0].Y = (gridHeight - 1) * scale
			}
			continue
		}

		if snake.body[0].X == vec.X && snake.body[0].Y == vec.Y {
			snake.isDead = true
		}

		hold.X = vec.X
		hold.Y = vec.Y
		vec.X = prev.X
		vec.Y = prev.Y
		prev.X = hold.X
		prev.Y = hold.Y
	}
}

func draw() {
	if !snake.isDead {
		rect.X = apple.pos.X
		rect.Y = apple.pos.Y
		rl.DrawRectanglePro(rect, origin, 0, rl.Orange)
	}
	for i, vec := range snake.body {
		rect.X = vec.X
		rect.Y = vec.Y
		if i == 0 {
			continue
		}
		rl.DrawRectanglePro(rect, origin, rotation, rl.Lime)
	}
	rect.X = snake.body[0].X
	rect.Y = snake.body[0].Y
	rl.DrawRectanglePro(rect, origin, rotation, rl.Green)

	rl.DrawRectangle(0, 0, gridWidth*scale, scale, rl.DarkGray)
	rl.DrawText(fmt.Sprintf("Score: %d", (len(snake.body)-2)*scoreMultiplier), scale/2, 1, scale, rl.White)

	if snake.isDead {
		rl.DrawText(fmt.Sprintf("Final Score: %d", (len(snake.body)-2)*scoreMultiplier), scale/2, scale*2, scale*2, rl.White)
	}

	if showFPS {
		rl.DrawFPS(scale/2, screenHeight-scale)
	}
}
