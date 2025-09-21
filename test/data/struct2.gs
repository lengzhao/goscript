// Interface example script
package main

type Rectangle struct {
	width  int
	height int
}

func (r Rectangle) Area() int {
	return r.width * r.height
}

// 没带指针的方法，修改应该是无效的
func (r Rectangle)SetWidth(width int) {
	r.width = width
}

func (r *Rectangle)SetHeight(height int) {
	r.height = height
}


func main() {
	// Create shapes
	rect := Rectangle{width: 10, height: 5}

	// 修改应该是无效的
	rect.SetWidth(6)
	// 修改应该是有效的
	rect.SetHeight(6)
	
	// Calculate areas
	rectArea := rect.Area()
	
	return rectArea
}