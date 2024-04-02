package main

import (
	"fmt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type rect2 struct {
	width, height float64
}

type circle struct {
	radius float64
}

// 要在 Go 中实现一个接口，我们就需要实现接口中的所有方法。
// 这里我们在 `rect` 上实现了 `geometry` 接口。
func (r rect2) area() float64 {
	return r.width * r.height
}

func (r rect2) perim() float64 {
	return 2*r.width + 2*r.height
}

// `circle` 的实现。
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := rect2{width: 2, height: 4}
	cicle := circle{radius: 5}

	measure(r)
	measure(cicle)
}
