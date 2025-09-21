// While loop script (using for loop with only condition)
package main

func countdown(n int) int {
    result := 0
    for n > 0 {
        result += n
        n--
    }
    return result
}

func main() {
    result := countdown(5)  // 5 + 4 + 3 + 2 + 1 = 15
    return result
}