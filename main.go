package main

import (
	"github.com/davecgh/go-spew/spew"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	n := Quad[Point]{
		Rectangle: Rectangle{
			Point: Point{
				x: 0,
				y: 0,
			},
			w: 64,
			h: 64,
		},
		IsLeaf:    true,
		Entities:  []Point{},
		SizeLimit: 2,
	}
	n.Insert(Point{x: 1, y: 1})
	n.Insert(Point{x: 1, y: 2})
	n.Insert(Point{x: 3, y: 3})
	err := n.Insert(Point{x: 63, y: 63})
	if err != nil {
		panic(err)
	}

	spew.Dump(n.Search(n))

	rl.InitWindow(640, 640, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var cellSize int = 10

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		for _, ent := range n.Search(n.Rectangle) {
			rl.DrawRectangle(int32(ent.x*cellSize), int32(ent.y*cellSize), int32(cellSize), int32(cellSize), rl.Red)
		}

		rl.EndDrawing()
	}
}
