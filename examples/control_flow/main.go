package main

import "fmt"

func main() {
	// Test various control flow constructs

	// Test if statement
	x := 10
	if x > 5 {
		fmt.Println("x is greater than 5")
	} else {
		fmt.Println("x is not greater than 5")
	}

	// Test if without else
	if x < 20 {
		fmt.Println("x is less than 20")
	}

	// Test switch statement
	switch x {
	case 5:
		fmt.Println("x is 5")
	case 10:
		fmt.Println("x is 10")
	default:
		fmt.Println("x is something else")
	}

	// Test for loop
	sum := 0
	for i := 1; i <= 3; i++ {
		sum += i
	}
	fmt.Printf("Sum from 1 to 3 is %d\n", sum)

	// Test range loop
	numbers := []int{1, 2, 3}
	rangeSum := 0
	for _, num := range numbers {
		rangeSum += num
	}
	fmt.Printf("Sum of numbers array is %d\n", rangeSum)

	// Test goto
	i := 0
	goto skip
	i = 100 // This should be skipped
skip:
	fmt.Printf("i is %d\n", i)

	fmt.Println("All control flow tests completed!")
}
