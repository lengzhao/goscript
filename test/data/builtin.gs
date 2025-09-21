// Test script for builtin functions
package main

func main() {
	// Test len function
	str := "hello"
	strLen := len(str)
	
	// Create a slice
	slice := make("slice", 5)
	sliceLen := len(slice)
	
	// Test copy function
	src := make("slice", 3)
	src[0] = 1
	src[1] = 2
	src[2] = 3
	
	dst := make("slice", 3)
	copied := copy(dst, src)
	
	// Print results
	print("String length:", strLen)
	print("Slice length:", sliceLen)
	print("Copied elements:", copied)
	print("Destination slice:", dst[0], dst[1], dst[2])
	
	return strLen + sliceLen + copied
}