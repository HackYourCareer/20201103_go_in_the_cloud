package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	r float64
}

type Rectangle struct {
	a float64
	b float64
}

func (c Circle) Area() float64 {
	return 2 * math.Pi * c.r
}

func (r Rectangle) Area() float64 {
	return r.a * r.b
}

func printArea(s Shape) {
	fmt.Printf("Shape's area: %v \n", s.Area())
}
