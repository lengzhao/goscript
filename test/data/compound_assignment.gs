// Compound assignment script
package main

func compoundOps() int {
    x := 10
    x += 5   // x = 15
    x -= 3   // x = 12
    x *= 2   // x = 24
    x /= 4   // x = 6
    return x
}

func main() {
    result := compoundOps()  // Should return 6
    return result
}