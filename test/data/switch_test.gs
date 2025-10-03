package main

func main() {
	x := 2
	result := 0
	
	switch x {
	case 1:
		result = 10
	case 2:
		result = 20
	case 3:
		result = 30
	default:
		result = 0
	}
	
	return result
}