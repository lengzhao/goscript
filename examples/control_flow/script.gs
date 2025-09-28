package main

func main() {
	x := 10
	y := 20

	// Test if statement
	if x < y {
		x = x + 5
	} else {
		y = y + 5
	}

	// Test for loop
	sum := 0
	for i := 0; i < 3; i++ {
		sum = sum + i
	}

	// Return the result
	return x + y + sum
}
