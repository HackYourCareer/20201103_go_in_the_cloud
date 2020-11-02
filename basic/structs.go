package main

import (
	"math"
)

type Complex struct {
	Im float64
	Re float64
}

func (c Complex) Modulus() float64 {
	sum := math.Pow(c.Im, 2) + math.Pow(c.Re, 2)

	return math.Sqrt(sum)
}
