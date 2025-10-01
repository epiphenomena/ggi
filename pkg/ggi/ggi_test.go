package ggi

import (
	"os"
	"testing"
)

func TestIsCGI(t *testing.T) {
	// Save original environment
	origMethod := os.Getenv("REQUEST_METHOD")
	origScript := os.Getenv("SCRIPT_NAME")
	
	// Test CGI detection
	os.Setenv("REQUEST_METHOD", "POST")
	os.Setenv("SCRIPT_NAME", "/test.cgi")
	
	if !IsCGI() {
		t.Error("Expected IsCGI() to return true when REQUEST_METHOD and SCRIPT_NAME are set")
	}
	
	// Reset environment
	os.Unsetenv("REQUEST_METHOD")
	os.Unsetenv("SCRIPT_NAME")
	
	if IsCGI() {
		t.Error("Expected IsCGI() to return false when REQUEST_METHOD and SCRIPT_NAME are not set")
	}
	
	// Restore original values
	if origMethod != "" {
		os.Setenv("REQUEST_METHOD", origMethod)
	} else {
		os.Unsetenv("REQUEST_METHOD")
	}
	if origScript != "" {
		os.Setenv("SCRIPT_NAME", origScript)
	} else {
		os.Unsetenv("SCRIPT_NAME")
	}
}