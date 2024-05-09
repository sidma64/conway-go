package main

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	board   Board
	running bool
	step    int
}
type Board [][]bool

const gameX = 100
const gameY = 100

type Color struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	return c.R, c.G, c.B, c.A
}

var cellClr = Color{0xFFFF, 0xFFFF, 0xFFFF, 0x0000}
var bgClr = Color{0x0000, 0x0000, 0x0000, 0x0000}

func NewBoard(x int, y int) (b Board, err error) {
	if x < 1 || y < 1 {
		return b, errors.New("You can't make a game with x or y coordinate that is lower than 1")
	}
	b = make([][]bool, x)
	for i := range x {
		b[i] = make([]bool, y)
	}
	return b, nil
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.running = true
	}
	if x, y := ebiten.CursorPosition(); ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		g.board.setCell(true, x, y)
	}
	if g.running {
		g.step++
		if g.step >= ebiten.TPS()/60 {
			g.step = 0
			g.board = g.board.next()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgClr)
	for i, row := range g.board {
		for j, alive := range row {
			if alive {
				screen.Set(i, j, cellClr)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gameX, gameY
}

func main() {
	b, _ := NewBoard(gameX, gameY)
	g := Game{
		board: b,
	}
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func (b Board) neighborCount(x int, y int) int {
	var count int
	for i := x - 1; i < x+2; i++ {
		if i < 0 || i >= len(b) {
			continue
		}
		for j := y - 1; j < y+2; j++ {
			if j < 0 || j >= len(b[0]) {
				continue
			}
			if i == x && y == j {
				continue
			}
			if b[i][j] {
				count++
			}
		}
	}
	return count
}

func (b Board) updateCell(x int, y int) bool {
	nc := b.neighborCount(x, y)
	alive := b[x][y]
	if !alive {
		if nc == 3 {
			return true
		}
		return false
	}
	if nc < 2 {
		return false
	}
	if nc == 2 || nc == 3 {
		return true
	}
	return false
}

func (b Board) next() Board {
	nextB, _ := NewBoard(len(b), len(b[0]))
	for i, col := range b {
		for j := range col {
			nextB[i][j] = b.updateCell(i, j)
		}
	}
	return nextB
}

func (b *Board) setCell(val bool, x int, y int) {
	if x < 0 || y < 0 || x >= gameX || y >= gameY {
		return
	}
	(*b)[x][y] = val
}
