package main

import (
	"fmt"
	"log"

	goscript "github.com/lengzhao/goscript"
)

func main() {
	// Example 1: Using strings module
	fmt.Println("=== Example 1: Using strings module ===")
	example1()

	// Example 2: Using json module
	fmt.Println("\n=== Example 2: Using json module ===")
	example2()

	// Example 3: Using modules via API
	fmt.Println("\n=== Example 3: Using modules via API ===")
	example3()
}

func example1() {
	// Create a new script
	script := goscript.NewScript([]byte(""))

	// Call strings module functions directly
	result1, err := script.CallFunction("strings.toLower", "HELLO WORLD")
	if err != nil {
		log.Printf("Error calling strings.toLower: %v", err)
	} else {
		fmt.Printf("strings.toLower('HELLO WORLD'): %v\n", result1)
	}

	result2, err := script.CallFunction("strings.contains", "hello world", "world")
	if err != nil {
		log.Printf("Error calling strings.contains: %v", err)
	} else {
		fmt.Printf("strings.contains('hello world', 'world'): %v\n", result2)
	}
}

func example2() {
	// Create a new script
	script := goscript.NewScript([]byte(""))

	// Call json module functions directly
	testMap := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	result1, err := script.CallFunction("json.marshal", testMap)
	if err != nil {
		log.Printf("Error calling json.marshal: %v", err)
	} else {
		fmt.Printf("json.marshal(map): %v\n", result1)
	}

	jsonStr := `{"name":"John","age":30}`
	result2, err := script.CallFunction("json.unmarshal", jsonStr)
	if err != nil {
		log.Printf("Error calling json.unmarshal: %v", err)
	} else {
		fmt.Printf("json.unmarshal(string): %v\n", result2)
	}
}

func example3() {
	// Example of using modules via API
	// Create a new script
	script := goscript.NewScript([]byte(""))

	// Import modules via API
	err := script.ImportModule("strings")
	if err != nil {
		log.Printf("Error importing strings module: %v", err)
		return
	}

	err = script.ImportModule("json")
	if err != nil {
		log.Printf("Error importing json module: %v", err)
		return
	}

	// Now we can call module functions directly
	fmt.Println("Calling module functions after importing:")

	lowerResult, err := script.CallFunction("strings.toLower", "HELLO WORLD")
	if err != nil {
		log.Printf("Error calling strings.toLower: %v", err)
	} else {
		fmt.Printf("strings.toLower('HELLO WORLD'): %v\n", lowerResult)
	}

	containsResult, err := script.CallFunction("strings.contains", "hello world", "world")
	if err != nil {
		log.Printf("Error calling strings.contains: %v", err)
	} else {
		fmt.Printf("strings.contains('hello world', 'world'): %v\n", containsResult)
	}

	// Test json functions
	testMap := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	marshalResult, err := script.CallFunction("json.marshal", testMap)
	if err != nil {
		log.Printf("Error calling json.marshal: %v", err)
	} else {
		fmt.Printf("json.marshal(map): %v\n", marshalResult)
	}

	jsonStr := `{"name":"John","age":30}`
	unmarshalResult, err := script.CallFunction("json.unmarshal", jsonStr)
	if err != nil {
		log.Printf("Error calling json.unmarshal: %v", err)
	} else {
		fmt.Printf("json.unmarshal(string): %v\n", unmarshalResult)
	}
}
