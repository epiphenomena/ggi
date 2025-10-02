package main

import (
	"fmt"
	"log"

	"ggi/pkg/ggi"
)

// Define custom content structures for our specific site
type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type PageContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	// Create an admin server to handle content management
	adminServer := ggi.NewAdminServer()
	
	// Register content types that this specific site needs
	// For example, register a markdown content for the homepage
	homeContent, err := ggi.NewMarkdownContent("_source/markdown/home.md")
	if err != nil {
		log.Fatal("Error creating markdown content:", err)
	}
	adminServer.RegisterContent(homeContent)
	
	// Register JSON content for cards
	cardsContent, err := ggi.NewJSONContent("_source/data/cards.json", []Card{})
	if err != nil {
		log.Fatal("Error creating JSON content:", err)
	}
	adminServer.RegisterContent(cardsContent)
	
	// Register media content
	logoContent, err := ggi.NewMediaContent("_source/media/logo.png")
	if err != nil {
		log.Fatal("Error creating media content:", err)
	}
	adminServer.RegisterContent(logoContent)
	
	// Build the site based on configuration
	config := ggi.BuildConfig{
		PublicTemplatesDir: "_templates/public", // Templates for public pages
		AdminTemplatesDir:  "_templates/admin",  // Templates for admin UI (not used in this basic example)
		ContentDir:         "_source",           // Directory for editable content
		OutputDir:          "_output",           // Directory for generated site
		EnableAdmin:        true,                // Generate admin UI
		EnableCGI:          true,                // Generate CGI script
	}

	fmt.Println("Building customized site...")
	
	err = ggi.Build(config)
	if err != nil {
		log.Fatal("Build failed:", err)
	}
	
	fmt.Println("Build completed successfully!")
	fmt.Println("Generated files:")
	fmt.Println("- _output/: Public static site")
	fmt.Println("- _output/admin/: Admin interface")
	fmt.Println("- _output/admin.cgi: CGI script for handling updates")
	fmt.Println("- _source/: Editable content files")
	
	fmt.Println("\nThe generated site can now be uploaded to your web server.")
	fmt.Println("Remember to set up .htaccess for security (basic auth for admin section).")
	
	// Note: The adminServer would be used when running the CGI script
	// The server handles both CGI mode (when called by web server) and 
	// development mode (when run directly)
}