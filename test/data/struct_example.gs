// Struct example script
package main

type Point struct {
	x int
	y int
}

type Rectangle struct {
	topLeft     Point
	bottomRight Point
}

func main() {
	// Create points
	p1 := Point{x: 10, y: 20}
	p2 := Point{x: 30, y: 40}
	
	// Create rectangle
	rect := Rectangle{topLeft: p1, bottomRight: p2}
	
	// Access fields
	width := rect.bottomRight.x - rect.topLeft.x
	height := rect.bottomRight.y - rect.topLeft.y
	area := width * height
	
	return area
}