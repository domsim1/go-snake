package main

import (
	"github.com/domsim1/go-snake/internal"
	"github.com/domsim1/go-snake/internal/flags"
	"github.com/domsim1/go-snake/pkg/scene"
	"github.com/domsim1/go-snake/pkg/scenes/snake"
	"github.com/domsim1/go-snake/pkg/scenes/title"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var sm = scene.NewSceneManager()

func main() {
	rl.InitWindow(internal.ScreenWidth, internal.ScreenHeight, "Go Snake")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	sm.Add("snake", snake.NewSnakeScene(sm))
	sm.Add("title", title.NewTitleScene(sm))
	sm.Activate("title")
	for !rl.WindowShouldClose() {
		update()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		draw()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func update() {
	if rl.IsKeyPressed(rl.KeyV) {
		if flags.VSync {
			flags.VSync = !flags.VSync
			rl.SetTargetFPS(0)
		} else {
			flags.VSync = !flags.VSync
			rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
		}
	}
	if rl.IsKeyPressed(rl.KeyF) {
		flags.ShowFPS = !flags.ShowFPS
	}
	sm.Active().Update()
}

func draw() {
	sm.Active().Draw()
	if flags.ShowFPS {
		rl.DrawFPS(internal.Scale/2, internal.ScreenHeight-internal.Scale)
	}
}
