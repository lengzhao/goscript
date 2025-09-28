package mypackage

// Global variable
var globalVar = 42

// Function with global variable access
func GetGlobalVar() int {
	return globalVar
}

// Function with local variable
func Calculate(x, y int) int {
	localVar := x + y
	return localVar
}

func main() {
	result := Calculate(10, 20)
	globalVar = globalVar + result
	return globalVar
}
