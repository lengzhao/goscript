package test

import (
	"fmt"
	"testing"
	"time"

	goscript "github.com/lengzhao/goscript"
)

// fibonacci calculates the nth Fibonacci number using pure Go
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// fibonacciIterative calculates the nth Fibonacci number using iterative approach in pure Go
func fibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// TestFibonacciPerformance compares the performance of pure Go vs GoScript implementations
func TestFibonacciPerformance(t *testing.T) {
	// t.Skip("Skipping performance test")
	// Test cases for different Fibonacci numbers
	testCases := []int{10, 20}

	for _, n := range testCases {
		t.Run(fmt.Sprintf("Fibonacci_%d", n), func(t *testing.T) {
			// Test pure Go recursive implementation
			start := time.Now()
			goResult := fibonacci(n)
			goDuration := time.Since(start)

			// Test pure Go iterative implementation
			start = time.Now()
			goIterResult := fibonacciIterative(n)
			goIterDuration := time.Since(start)

			// Test GoScript implementation
			scriptSource := fmt.Sprintf(`
package main

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	result := fibonacci(%d)
	return result
}
`, n)

			script := goscript.NewScript([]byte(scriptSource))
			script.SetDebug(false) // Disable debug to reduce output
			start = time.Now()
			scriptResult, err := script.Run()
			scriptDuration := time.Since(start)

			if err != nil {
				t.Fatalf("Failed to run GoScript: %v", err)
			}

			// Verify results match
			if goResult != scriptResult {
				t.Errorf("Results don't match: Go=%d, GoScript=%d", goResult, scriptResult)
			}

			if goResult != goIterResult {
				t.Errorf("Results don't match: Go=%d, GoIterative=%d", goResult, goIterResult)
			}

			// Log performance results
			t.Logf("Fibonacci(%d) Results:", n)
			t.Logf("  Pure Go (recursive): %d, Time: %v", goResult, goDuration)
			t.Logf("  Pure Go (iterative): %d, Time: %v", goIterResult, goIterDuration)
			t.Logf("  GoScript: %d, Time: %v", scriptResult, scriptDuration)
			t.Logf("  Speedup (GoScript/Pure Go recursive): %.2fx", float64(scriptDuration)/float64(goDuration))
			t.Logf("  Speedup (GoScript/Pure Go iterative): %.2fx", float64(scriptDuration)/float64(goIterDuration))
		})
	}
	// t.Error("Benchmarking completed")
}

// BenchmarkFibonacci benchmarks the pure Go implementation
func BenchmarkFibonacci(b *testing.B) {
	n := 20
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fibonacci(n)
	}
}

// BenchmarkFibonacciIterative benchmarks the pure Go iterative implementation
func BenchmarkFibonacciIterative(b *testing.B) {
	n := 20
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fibonacciIterative(n)
	}
}

// BenchmarkFibonacciScript benchmarks the GoScript implementation
func BenchmarkFibonacciScript(b *testing.B) {
	scriptSource := `
package main

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	result := fibonacci(20)
	return result
}
`
	script := goscript.NewScript([]byte(scriptSource))
	script.SetDebug(false) // Disable debug to reduce output

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := script.Run()
		if err != nil {
			b.Fatalf("Failed to run GoScript: %v", err)
		}
	}
}
