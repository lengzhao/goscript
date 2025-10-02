package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/lengzhao/goscript"
)

func main() {
	n := flag.Int("n", 30, "The number to calculate fibonacci")
	profile := flag.Bool("profile", false, "Enable CPU profiling")
	flag.Parse()

	// Start CPU profiling if requested
	if *profile {
		fmt.Println("Starting CPU profiling...")
		f, err := os.Create("fibonacci.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
`, *n)

	script := goscript.NewScript([]byte(scriptSource))
	script.SetDebug(false) // Disable debug to reduce output
	start := time.Now()
	scriptResult, err := script.Run()
	scriptDuration := time.Since(start)

	if err != nil {
		panic(err)
	}

	fmt.Println("GoScript result:", scriptResult)
	fmt.Println("Script duration:", scriptDuration)

	// Print execution stats
	stats := script.GetExecutionStats()
	fmt.Printf("Execution time: %v\n", stats.ExecutionTime)
	fmt.Printf("Instruction count: %d\n", stats.InstructionCount)
	fmt.Printf("Error count: %d\n", stats.ErrorCount)
}
