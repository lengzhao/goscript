package main

func main() {
	// Use if-else statement
	x := 10
	if x > 5 {
		x = x * 2
	} else {
		x = x / 2
	}

	// Use for loop
	sum := 0
	for i := 1; i <= 3; i++ {  // Reduced iterations
		sum += i
	}

	// Use while-like loop
	count := 0
	for count < 2 {  // Reduced iterations
		count++
	}

	return x + sum + count
}