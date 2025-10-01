package ggi

import (
	"net/http"
	"os"
	"strings"
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

func TestAuthentication(t *testing.T) {
	// Save original secret key
	origKey := SecretKey
	
	// Set test key
	SecretKey = "test_secret_key"
	
	// Create a test request with correct token
	formData := "token=test_secret_key"
	req, err := http.NewRequest("POST", "/", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	if !IsAuthenticated(req) {
		t.Error("Expected IsAuthenticated to return true with correct token")
	}
	
	// Create a test request with wrong token
	formData = "token=wrong_token"
	req, err = http.NewRequest("POST", "/", strings.NewReader(formData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	if IsAuthenticated(req) {
		t.Error("Expected IsAuthenticated to return false with wrong token")
	}
	
	// Create a GET request (should fail authentication)
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	if IsAuthenticated(req) {
		t.Error("Expected IsAuthenticated to return false with GET request")
	}
	
	// Restore original key
	SecretKey = origKey
}