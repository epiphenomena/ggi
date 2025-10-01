package main

import (
	"fmt"
	"log"
	"os"

	"ggi/pkg/ggi"
)

func main() {
	// Set secret key as required by the ggi package
	if secret := os.Getenv("GGI_SECRET_KEY"); secret != "" {
		ggi.SecretKey = secret
	} else {
		// Use a default key for the build process
		ggi.SecretKey = "build_secret_key"
	}

	// Example of using the build functionality
	config := ggi.BuildConfig{
		SourceDir: "examples/basic-site/_source",
		OutputDir: "examples/basic-site/_output", 
		BaseURL:   "",
	}

	fmt.Println("Starting build process...")
	
	err := ggi.Build(config)
	if err != nil {
		log.Fatal("Build failed:", err)
	}
	
	fmt.Println("Build completed successfully!")
	fmt.Println("Static site generated in:", config.OutputDir)
	
	fmt.Println("You can now serve the files in the output directory with a web server.")
}