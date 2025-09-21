// Package main provides an example of using interfaces in GoScript
package main

import (
	"fmt"

	"github.com/lengzhao/goscript"
)

func main() {
	// Create a script that uses interfaces
	source := `
package main

type Shape interface {
	Area() float64
}

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
	rect := Rectangle{width: 10, height: 5}
	circle := Circle{radius: 3}

	// Calculate areas
	rectArea := rect.Area()
	circleArea := circle.Area()

	return rectArea + circleArea
}
`

	script := goscript.NewScript([]byte(source))

	// Execute the script
	result, err := script.Run()
	if err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	fmt.Printf("Script result: %v\n", result)
}
