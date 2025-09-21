// Loop script
package main

func sum(n int) int {
    total := 0
    for i := 1; i <= n; i++ {
        total += i
    }
    return total
}

func main() {
    result := sum(10)  // 1+2+3+...+10 = 55
    return result
}