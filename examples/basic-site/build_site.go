package main

import (
	"fmt"
	"log"

	"ggi/pkg/ggi"
)

func main() {
	// Example of using the build functionality with separate public and admin sections
	config := ggi.BuildConfig{
		PublicSourceDir: "examples/basic-site/_source/public",
		AdminSourceDir:  "examples/basic-site/_source/admin",
		OutputDir:       "examples/basic-site/_output", 
		BaseURL:         "",
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