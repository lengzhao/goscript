// Complex test script for OpRegistFunction
package main

func add(a, b int) int {
    return a + b
}

func multiply(a, b int) int {
    return a * b
}

func calculate(x, y int) int {
    sum := add(x, y)
    product := multiply(x, y)
    return sum + product
}

func main() {
    result1 := add(3, 4)
    result2 := multiply(5, 6)
    result3 := calculate(2, 3)
    return result1 + result2 + result3
}