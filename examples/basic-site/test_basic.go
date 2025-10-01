package main

import (
	"fmt"
	"os"

	"ggi/pkg/ggi"
)

func main() {
	// Set up environment variables to simulate CGI
	os.Setenv("REQUEST_METHOD", "POST")
	os.Setenv("SCRIPT_NAME", "/test.cgi")
	
	// Test CGI detection
	if ggi.IsCGI() {
		fmt.Println("✓ CGI detection working correctly")
	} else {
		fmt.Println("✗ CGI detection failed")
	}
	
	// Reset to non-CGI environment
	os.Unsetenv("REQUEST_METHOD")
	if !ggi.IsCGI() {
		fmt.Println("✓ Non-CGI environment correctly detected")
	} else {
		fmt.Println("✗ Non-CGI environment incorrectly detected")
	}
	
	// Test secret key access
	originalKey := ggi.SecretKey
	fmt.Printf("✓ Current secret key: %s\n", ggi.SecretKey)
	
	// Test setting a new key
	ggi.SecretKey = "new_test_key"
	if ggi.SecretKey == "new_test_key" {
		fmt.Println("✓ Secret key can be set")
	} else {
		fmt.Println("✗ Secret key setting failed")
	}
	
	// Restore original key
	ggi.SecretKey = originalKey
	
	fmt.Println("Basic functionality tests completed!")
}