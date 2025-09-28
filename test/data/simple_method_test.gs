package main

type Rectangle struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (c Circle) Area() float64 {
	return 3.14159 * c.radius * c.radius
}

func main() {
	// Create shapes
	rect := Rectangle{width: 10.0, height: 5.0}
	circle := Circle{radius: 3.0}
	
	// Calculate areas
	rectArea := rect.Area()
	circleArea := circle.Area()
	
	return rectArea + circleArea
}