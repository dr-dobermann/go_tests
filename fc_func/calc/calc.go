package main

import (
	"fmt"
	"math"
)

type Calculator struct {
	acc float64
}

type opfunc func(float64) float64

func (c *Calculator) Do(op opfunc) float64 {
	c.acc = op(c.acc)
	return c.acc
}

func Add(n float64) opfunc {
	return func(acc float64) float64 {
		return acc + n
	}
}

func Sub(n float64) opfunc {
	return func(acc float64) float64 {
		return acc - n
	}
}

func Mul(n float64) opfunc {
	return func(acc float64) float64 {
		return acc * n
	}
}

func Sqrt() opfunc {
	return func(n float64) float64 {
		return math.Sqrt(n)
	}
}

func main() {
	var c Calculator
	fmt.Println(c.Do(Add(5)))
	fmt.Println(c.Do(Sub(3)))
	fmt.Println(c.Do(Mul(8)))
	fmt.Println(c.Do(Sqrt()))
	fmt.Println(c.Do(math.Sin))
}
