package snake

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/domsim1/go-snake/internal"
	"github.com/domsim1/go-snake/pkg/direction"
	"github.com/domsim1/go-snake/pkg/scene"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	startingTickRate = 0.25
	appleTickLimit   = (internal.GridWidth + internal.GridHeight) * 0.8
	speedLimit       = 0.060
	speedRate        = 0.005
	scoreMultiplier  = 10
)

type Apple struct {
	pos    *rl.Vector2
	ticker int
}

type Snake struct {
	body          []*rl.Vector2
	timer         float32
	tickRate      float32
	isDead        bool
	direction     direction.Direction
	nextDirection direction.Direction
}

type state struct {
	isPaused bool
	snake    *Snake
	rect     *rl.Rectangle
	rotation float32
	origin   *rl.Vector2
	apple    *Apple
	sm       scene.SceneManager
	eatSound rl.Sound
	isMuted  bool
}

func NewSnakeScene(sm scene.SceneManager) scene.Scene {
	ex, err := os.Executable()
	if err != nil {
		return nil
	}
	exPath := filepath.Dir(ex)

	return &state{
		origin: &rl.Vector2{
			X: 0,
			Y: 0,
		},
		sm: sm,
		eatSound: rl.LoadSound(exPath + "/resources/eat.ogg"),
		isMuted: false,
	}
}

func (s *state) Setup() {
	s.snake = &Snake{
		body:          make([]*rl.Vector2, 0),
		timer:         0.0,
		tickRate:      startingTickRate,
		isDead:        false,
		direction:     direction.Right,
		nextDirection: direction.Right,
	}
	s.snake.body = append(s.snake.body, &rl.Vector2{
		X: (internal.GridWidth / 2) * internal.Scale,
		Y: (internal.GridHeight / 2) * internal.Scale,
	})
	s.snake.body = append(s.snake.body, &rl.Vector2{
		X: ((internal.GridWidth / 2) - 1) * internal.Scale,
		Y: (internal.GridHeight / 2) * internal.Scale,
	})
	s.rect = &rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  internal.Scale,
		Height: internal.Scale,
	}
	s.apple = &Apple{
		pos: &rl.Vector2{
			X: float32(rl.GetRandomValue(1, internal.GridWidth-1) * internal.Scale),
			Y: float32(rl.GetRandomValue(1, internal.GridHeight-1) * internal.Scale),
		},
		ticker: 0,
	}
	s.rotation = 0.0
	s.isPaused = false
	s.isMuted = true
}

func (s *state) Update() {	
	if rl.IsKeyPressed(rl.KeyR) {
		s.Setup()
		return
	}
	if !rl.IsWindowFocused() && !s.snake.isDead {
		s.isPaused = true
	}
	if rl.IsKeyPressed(rl.KeyP) && !s.snake.isDead && rl.IsWindowFocused() {
		s.isPaused = !s.isPaused
	}
	if rl.IsKeyPressed(rl.KeyM) {
		s.isMuted = !s.isMuted
	}
	if s.isPaused && !s.snake.isDead {
		return
	}
	if rl.IsKeyDown(rl.KeyD) && s.snake.direction != direction.Left {
		s.snake.nextDirection = direction.Right
	}
	if rl.IsKeyDown(rl.KeyA) && s.snake.direction != direction.Right {
		s.snake.nextDirection = direction.Left
	}
	if rl.IsKeyDown(rl.KeyW) && s.snake.direction != direction.Down {
		s.snake.nextDirection = direction.Up
	}
	if rl.IsKeyDown(rl.KeyS) && s.snake.direction != direction.Up {
		s.snake.nextDirection = direction.Down
	}
	dt := rl.GetFrameTime()
	if s.snake.isDead {
		s.rotation += 60.0 * dt
		if s.rotation > 360 {
			s.rotation = 0
		}
		return
	}

	s.snake.timer += dt
	if s.snake.timer < s.snake.tickRate {
		return
	}
	s.snake.timer = 0.0
	s.snake.direction = s.snake.nextDirection

	if s.snake.body[0].X == s.apple.pos.X && s.snake.body[0].Y == s.apple.pos.Y {
		if !s.isMuted {
			rl.PlaySound(s.eatSound)
		}

		s.snake.body = append(s.snake.body, &rl.Vector2{
			X: s.snake.body[len(s.snake.body)-1].X,
			Y: s.snake.body[len(s.snake.body)-1].Y,
		})
		s.apple.pos.X = float32(rl.GetRandomValue(1, internal.GridWidth-1) * internal.Scale)
		s.apple.pos.Y = float32(rl.GetRandomValue(1, internal.GridHeight-1) * internal.Scale)
		s.apple.ticker = 0
		if s.snake.tickRate > speedLimit {
			s.snake.tickRate -= speedRate
		}
	}

	s.apple.ticker += 1
	if s.apple.ticker > appleTickLimit {
		s.apple.pos.X = float32(rl.GetRandomValue(1, internal.GridWidth-1) * internal.Scale)
		s.apple.pos.Y = float32(rl.GetRandomValue(1, internal.GridHeight-1) * internal.Scale)
		s.apple.ticker = 0
	}

	prev := rl.Vector2{
		X: s.snake.body[0].X,
		Y: s.snake.body[0].Y,
	}
	hold := rl.Vector2{
		X: s.snake.body[0].X,
		Y: s.snake.body[0].Y,
	}
	for i, vec := range s.snake.body {
		if i == 0 {
			switch s.snake.direction {
			case direction.Up:
				vec.Y -= internal.Scale
			case direction.Down:
				vec.Y += internal.Scale
			case direction.Left:
				vec.X -= internal.Scale
			case direction.Right:
				vec.X += internal.Scale
			}

			if s.snake.body[0].X > (internal.GridWidth-1)*internal.Scale {
				s.snake.body[0].X = 0
			}
			if s.snake.body[0].X < 0 {
				s.snake.body[0].X = (internal.GridWidth - 1) * internal.Scale
			}
			if s.snake.body[0].Y > (internal.GridHeight-1)*internal.Scale {
				s.snake.body[0].Y = internal.Scale
			}
			if s.snake.body[0].Y < internal.Scale {
				s.snake.body[0].Y = (internal.GridHeight - 1) * internal.Scale
			}
			continue
		}

		if s.snake.body[0].X == vec.X && s.snake.body[0].Y == vec.Y {
			s.snake.isDead = true
		}

		hold.X = vec.X
		hold.Y = vec.Y
		vec.X = prev.X
		vec.Y = prev.Y
		prev.X = hold.X
		prev.Y = hold.Y
	}
}

func (s *state) Draw() {
	if !s.snake.isDead {
		s.rect.X = s.apple.pos.X
		s.rect.Y = s.apple.pos.Y
		rl.DrawRectanglePro(*s.rect, *s.origin, 0, rl.Orange)
	}
	for i, vec := range s.snake.body {
		s.rect.X = vec.X
		s.rect.Y = vec.Y
		if i == 0 {
			continue
		}
		rl.DrawRectanglePro(*s.rect, *s.origin, s.rotation, rl.Lime)
	}
	s.rect.X = s.snake.body[0].X
	s.rect.Y = s.snake.body[0].Y
	rl.DrawRectanglePro(*s.rect, *s.origin, s.rotation, rl.Green)

	rl.DrawRectangle(0, 0, internal.GridWidth*internal.Scale, internal.Scale, rl.DarkGray)
	rl.DrawText(fmt.Sprintf("Score: %d", (len(s.snake.body)-2)*scoreMultiplier), internal.Scale/2, 1, internal.Scale, rl.White)

	if s.snake.isDead {
		rl.DrawText(fmt.Sprintf("Final Score: %d", (len(s.snake.body)-2)*scoreMultiplier), internal.Scale/2, internal.Scale*2, internal.Scale*2, rl.White)
		rl.DrawText("press r to restart", internal.Scale/2, internal.Scale*5, internal.Scale, rl.White)
	}
	if s.isPaused {
		rl.DrawText("-- paused --", internal.Scale/2, internal.Scale*2, internal.Scale, rl.White)
	}

	if s.isMuted {
		rl.DrawText("-- muted --", internal.ScreenWidth - (8*13), 1, internal.Scale, rl.White)
	}
}
