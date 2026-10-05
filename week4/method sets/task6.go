package main

import (
	"fmt"
	"math"
)

// Task 6: Shapes
// 1. Declare type Shape interface { Area() float64; Perimeter() float64 } .
// 2. Implement Rectangle with value receivers, plus a method Scale(k
// float64) that really changes the rectangle.
// 3. Implement Circle with pointer receivers.
// 4. Put one rectangle and one circle into a []Shape and print the area and
// perimeter of each in a loop.
// 5. Try to put Circle{Radius: 1} (not &Circle{...} ) into the slice. Read the
// compiler error and explain it.

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Scale has to be a POINTER receiver, otherwise it would scale a copy and the rectangle wouldn't change.
// (This mixes receivers on Rectangle, but it's required by the task. Area/Perimeter stay value receivers
// so Rectangle{} itself still implements Shape.)
func (r *Rectangle) Scale(k float64) {
	r.Width *= k
	r.Height *= k
}

type Circle struct{ Radius float64 }

func (c *Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c *Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

func main() {
	rect := Rectangle{3, 4}
	rect.Scale(2)     // rect is a variable -> addressable -> (&rect).Scale(2)
	fmt.Println(rect) // {6 8}

	shapes := []Shape{
		rect,               // value is fine: Rectangle's method set has Area and Perimeter
		&Circle{Radius: 1}, // must be a pointer
	}
	for _, s := range shapes {
		fmt.Printf("%T area=%.2f perimeter=%.2f\n", s, s.Area(), s.Perimeter())
	}
	// main.Rectangle area=48.00 perimeter=28.00
	// *main.Circle area=3.14 perimeter=6.28

	// 5.
	// shapes = append(shapes, Circle{Radius: 1})
	// cannot use Circle{…} (value of struct type Circle) as Shape value in argument to append:
	// Circle does not implement Shape (method Area has pointer receiver)
	//
	// Area and Perimeter are declared on *Circle, so they are in the method set of *Circle only. The method set of
	// Circle (the value) is empty. The value inside an interface is a non-addressable copy, so Go can't
	// call a pointer method on it. Use &Circle{Radius: 1}.
}
