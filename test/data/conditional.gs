// Conditional script
package main

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main() {
    result := max(10, 20)
    return result
}