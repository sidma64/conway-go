package main

import (
	"errors"
)

type Point struct {
	x int
	y int
}

func (ent Point) X() int {
	return ent.x
}

func (ent Point) Y() int {
	return ent.y
}

func (ent Point) H() int {
	return 1
}

func (ent Point) W() int {
	return 1
}

type Searchable interface {
	X() int
	Y() int
	H() int
	W() int
}

type Rectangle struct {
	Point
	w int
	h int
}

func (r Rectangle) X() int {
	return r.x
}

func (r Rectangle) Y() int {
	return r.y
}

func (r Rectangle) H() int {
	return r.h
}

func (r Rectangle) W() int {
	return r.w
}

func Intersects(a Searchable, b Searchable) bool {
	return a.X() < b.X()+b.W() && a.X()+a.W() > b.X() && a.Y() < b.Y()+b.H() && a.Y()+a.H() > b.Y()
}

type Quad[T Searchable] struct {
	Rectangle
	Entities    []T
	IsLeaf      bool
	TopLeft     *Quad[T]
	TopRight    *Quad[T]
	BottomLeft  *Quad[T]
	BottomRight *Quad[T]
	SizeLimit   int
}

func (n *Quad[T]) Slice() {
	mid := n.w / 2
	top := Rectangle{
		Point: Point{
			x: n.x,
			y: n.y,
		},
		w: n.w,
		h: mid,
	}
	bot := Rectangle{
		Point: Point{
			x: n.x,
			y: n.y + mid,
		},
		w: n.w,
		h: mid,
	}

	var topLeft []T
	for _, v := range n.Entities {
		if Intersects(top, v) {
			topLeft = append(topLeft, v)
		}
	}
	var topRight []T
	for _, v := range n.Entities {
		if Intersects(top, v) {
			topRight = append(topRight, v)
		}
	}
	var botLeft []T
	for _, v := range n.Entities {
		if Intersects(bot, v) {
			botLeft = append(botLeft, v)
		}
	}
	var botRight []T
	for _, v := range n.Entities {
		if Intersects(bot, v) {
			botRight = append(botRight, v)
		}
	}
	n.TopLeft = &Quad[T]{
		Rectangle: Rectangle{
			Point: Point{
				x: n.x,
				y: n.y,
			},
			w: mid,
			h: mid,
		},
		IsLeaf:    true,
		Entities:  topLeft,
		SizeLimit: n.SizeLimit,
	}
	n.TopRight = &Quad[T]{
		Rectangle: Rectangle{
			Point: Point{
				x: n.x + mid,
				y: n.y,
			},
			w: mid,
			h: mid,
		},
		IsLeaf:    true,
		Entities:  topRight,
		SizeLimit: n.SizeLimit,
	}
	n.BottomLeft = &Quad[T]{
		Rectangle: Rectangle{
			Point: Point{
				x: n.x,
				y: n.y + mid,
			},
			w: mid,
			h: mid,
		},
		IsLeaf:    true,
		Entities:  botLeft,
		SizeLimit: n.SizeLimit,
	}
	n.BottomRight = &Quad[T]{
		Rectangle: Rectangle{
			Point: Point{
				x: n.x + mid,
				y: n.y + mid,
			},
			w: mid,
			h: mid,
		},
		IsLeaf:    true,
		Entities:  botRight,
		SizeLimit: n.SizeLimit,
	}
	n.IsLeaf = false
}

func (n *Quad[T]) Search(query Searchable) []T {
	var result []T
	if !Intersects(n, query) {
		return result
	}
	if n.IsLeaf {
		for _, v := range n.Entities {
			if Intersects(query, v) {
				result = append(result, v)
			}

		}
		return result
	}

	result = append(result, n.TopLeft.Search(query)...)
	result = append(result, n.TopRight.Search(query)...)
	result = append(result, n.BottomLeft.Search(query)...)
	result = append(result, n.BottomRight.Search(query)...)
	return result
}

func (n *Quad[T]) Insert(ent T) error {
	if !Intersects(n, ent) {
		return errors.New("entity out of bounds")
	}

	if n.IsLeaf {
		if len(n.Entities) > n.SizeLimit {
			n.Slice()

		} else {
			for _, v := range n.Entities {
				if Intersects(ent, v) {
					return errors.New("entity already exists")
				}
			}
			n.Entities = append(n.Entities, ent)
			return nil
		}
	}
	_ = n.TopLeft.Insert(ent)
	_ = n.TopRight.Insert(ent)
	_ = n.BottomLeft.Insert(ent)
	err := n.BottomRight.Insert(ent)
	return err
}
