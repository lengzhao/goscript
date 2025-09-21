// Complex script with multiple operations
package main

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    fact := factorial(5)  // 120
    fib := fibonacci(10)  // 55
    return fact + fib     // 175
}