package main

import (
	"ggi/pkg/ggi"
)

// This will test that we get a panic when SecretKey is not set before using authentication functions
func main() {
	// Temporarily save the original key
	originalKey := ggi.SecretKey
	
	// Set it to empty to test the panic behavior
	ggi.SecretKey = ""
	
	// Now trying to use IsAuthenticated should panic
	// We'll test by running a defer/recover
	defer func() {
		if r := recover(); r != nil {
			println("Got expected panic:", r.(string))
		}
	}()
	
	// This should cause a panic since SecretKey is empty
	ggi.IsAuthenticated(nil) // This will panic
	
	// Restore the original key if we somehow get here
	ggi.SecretKey = originalKey
}