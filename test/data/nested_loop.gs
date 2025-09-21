// Nested loop script
package main

func multiplyTables() int {
    result := 0
    for i := 1; i <= 3; i++ {
        for j := 1; j <= 3; j++ {
            result += i * j
        }
    }
    return result
}

func main() {
    result := multiplyTables()  // (1*1 + 1*2 + 1*3) + (2*1 + 2*2 + 2*3) + (3*1 + 3*2 + 3*3) = 36
    return result
}