package math

func add(a, b) {
	return a + b
}

func multiply(a, b) {
	return a * b
}

func square(x) {
	return multiply(x, x)
}

func fibonacci(n) {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}