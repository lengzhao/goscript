package main

func main() {
	i := 0
	goto label1
	
	// This code should be skipped
	i += 100
	
label1:
	// This code should be executed
	i += 42
	
	return i
}