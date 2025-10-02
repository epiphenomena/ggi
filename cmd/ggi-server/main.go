package main

import (
	"ggi/pkg/ggi"
)

func main() {
	server := ggi.NewAdminServer()
	
	// Register content types that this site will handle
	// In a real generator, the professional would register the specific
	// content types needed for their site
	
	// For this example, register some common content types
	defaultMarkdown, _ := ggi.NewMarkdownContent("_source/content/default.md")
	server.RegisterContent(defaultMarkdown)
	
	defaultJSON, _ := ggi.NewJSONContent("_source/data/default.json", []interface{}{})
	server.RegisterContent(defaultJSON)
	
	defaultMedia, _ := ggi.NewMediaContent("_source/media/default.jpg")
	server.RegisterContent(defaultMedia)

	// Run server in CGI mode if detected, otherwise start dev server
	server.RunServer()
}