package main

import (
	"errors"
	"time"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	x int
	y int
}

type Life struct {
	board [][]bool
	w     int
	h     int
}
var ErrOutOfBounds = errors.New("out of bounds")

func NewLife(w int, h int) Life {
	b := make([][]bool, h)
	for i := 0; i < len(b); i++ {
		b[i] = make([]bool, w)
	}
	return Life{b, w, h}
}

func (l Life) Step() Life {
	n := NewLife(l.w, l.h)
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			p := Point{x, y}
			c := l.NeighborCount(p)
			if cell, _ := l.Get(p); cell {
				n.Put(p, c == 2 || c == 3)
			} else {
				n.Put(p, c == 3)
			}
		}
	}
	return n
}

func (l *Life) FillRandom(count int) {
	for i := 0; i < count; i++ {
		x := rand.IntN(l.w)
		y := rand.IntN(l.h)
		l.Put(Point{x, y}, true)
	}
}

func (l Life) NeighborCount(p Point) int {
	var c int
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if cell, _ := l.Get(Point{p.x + i, p.y + j}); cell {
				c++
			}
		}
	}
	return c
}

func (l Life) Contains(p Point) bool {
	return p.x >= 0 && p.x < l.w && p.y >= 0 && p.y < l.h
}

func (l *Life) Put(p Point, val bool) error {
	if l.Contains(p) {
		l.board[p.y][p.x] = val
	}
	return ErrOutOfBounds
}

func (l Life) Get(p Point) (bool, error) {
	if l.Contains(p) {
		return l.board[p.y][p.x], nil
	}
	return false, ErrOutOfBounds
}

type Game struct {
	life Life
	zoom float64
	prevUpdateTime time.Time
	sinceStep time.Duration
	stepSize time.Duration
}

func (g *Game) Update() {
	dt := time.Since(g.prevUpdateTime)
	g.prevUpdateTime = time.Now()
	g.sinceStep += dt
	if g.sinceStep > g.stepSize {
		g.life = g.life.Step()
		g.sinceStep = g.sinceStep - g.stepSize
	}

}

func main() {
	var g Game
	g.life = NewLife(1000, 1000)
	g.zoom = 1
	g.stepSize = time.Second * 2
	g.prevUpdateTime = time.Now()
	g.sinceStep = 0
	g.life.FillRandom(200000)
	rl.InitWindow(1000, 1000, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)
	rl.GetFrameTime()

	for !rl.WindowShouldClose() {
		g.Update()
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawFPS(0, 0)
		for y := 0; y < g.life.h; y++ {
			for x := 0; x < g.life.w; x++ {
				if cell, _ := g.life.Get(Point{x, y}); cell {
					rl.DrawRectangle(int32(x)*int32(g.zoom), int32(y)*int32(g.zoom), int32(g.zoom), int32(g.zoom), rl.Black)
				}
			}
		}
		rl.EndDrawing()
	}
}
