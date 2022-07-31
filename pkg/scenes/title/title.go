package title

import (
	"github.com/domsim1/go-snake/internal"
	"github.com/domsim1/go-snake/pkg/scene"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	timerLimit = 1.8
	hiddenTime = 1.6
)

type state struct {
	sm       scene.SceneManager
	selected int
	timer    float32
}

func NewTitleScene(sm scene.SceneManager) scene.Scene {
	return &state{
		sm: sm,
	}
}

func (s *state) Setup() {
	s.selected = 0
	s.timer = 0.0
}

func (s *state) Update() {
	if rl.IsKeyPressed(rl.KeyR) {
		s.sm.Activate("snake")
	}
	dt := rl.GetFrameTime()
	s.timer += dt
	if s.timer > timerLimit {
		s.timer = 0
	}
}

func (s *state) Draw() {
	rl.DrawText("GO SNAKE", internal.Scale, internal.Scale, internal.Scale*3, rl.Green)
	if s.timer < hiddenTime {
		rl.DrawText("PRESS R TO START", internal.Scale, internal.Scale*4, internal.Scale*2, rl.Orange)
	}
	rl.DrawText("--- controlls ---", internal.Scale, internal.Scale*7, internal.Scale, rl.White)
	rl.DrawText("w    move up", internal.Scale, internal.Scale*8, internal.Scale, rl.White)
	rl.DrawText("s    move down", internal.Scale, internal.Scale*9, internal.Scale, rl.White)
	rl.DrawText("a    move left", internal.Scale, internal.Scale*10, internal.Scale, rl.White)
	rl.DrawText("d    move right", internal.Scale, internal.Scale*11, internal.Scale, rl.White)
	rl.DrawText("r    restart", internal.Scale, internal.Scale*12, internal.Scale, rl.White)
	rl.DrawText("p    toggle pause", internal.Scale, internal.Scale*13, internal.Scale, rl.White)
	rl.DrawText("v    toggle vsync", internal.Scale, internal.Scale*14, internal.Scale, rl.White)
	rl.DrawText("f    toggle fps", internal.Scale, internal.Scale*15, internal.Scale, rl.White)
}
