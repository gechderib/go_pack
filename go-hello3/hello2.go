package main

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}
type Rectangle struct {
	Width  float64
	Height float64
}

type Triangle struct {
	Base   float64
	Height float64
}

func (circle *Circle) Area() float64 {
	return 3.14 * circle.Radius * circle.Radius
}

func (circle *Circle) Perimeter() float64 {
	return 2 * 3.14 * circle.Radius
}

func (rect *Rectangle) Area() float64 {
	return rect.Width * rect.Height
}

func (rect *Rectangle) Perimeter() float64 {
	return 2 * (rect.Width + rect.Height)
}

func (tri *Triangle) Area() float64 {
	return 0.5 * tri.Base * tri.Height
}

func (tri *Triangle) Perimeter() float64 {
	// Assuming an equilateral triangle for simplicity
	return 3 * tri.Base
}

func calculateArea(shape Shape) float64 {
	return shape.Area()
}

func CalculatePerimeter(shape Shape) float64 {
	return shape.Perimeter()
}
