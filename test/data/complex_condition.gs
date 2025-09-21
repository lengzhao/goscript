// Complex condition script
package main

func complexCondition(a, b, c int) int {
    if a > 0 && b > 0 {
        if c > 0 {
            return a + b + c
        } else {
            return a + b
        }
    } else if a < 0 || b < 0 {
        return a * b
    } else {
        return 0
    }
}

func main() {
    result := complexCondition(3, 4, 5)  // 3 > 0 && 4 > 0 && 5 > 0, so return 3 + 4 + 5 = 12
    return result
}