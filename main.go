package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const devServerPort = "8082"

var adminServer *AdminServer

func main() {
	serve := flag.Bool("serve", false, "Run development server")
	fastcgi := flag.Bool("fastcgi", false, "Run FastCGI server")
	build := flag.Bool("build", false, "Build the static site")
	clean := flag.Bool("clean", false, "Clean the public folder of build artifacts")
	flag.Parse()

	// Initialize admin server
	adminServer = NewAdminServer()

	if *serve {
		fmt.Println("Starting development server...")
		startDevServer()
	} else if *fastcgi {
		fmt.Println("Starting FastCGI server...")
		startFastCGIServer()
	} else if *build {
		fmt.Println("Building static site...")
		if err := Build(); err != nil {
			fmt.Printf("Error building site: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Build completed successfully!")
	} else if *clean {
		fmt.Println("Cleaning public folder...")
		if err := Clean(); err != nil {
			fmt.Printf("Error cleaning: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Clean completed successfully!")
	} else {
		// Run as CGI script
		fmt.Println("Running as CGI script...")
		runCGI()
	}
}

// runCGI handles the CGI script execution
func runCGI() {
	// Get the request method and path
	method := os.Getenv("REQUEST_METHOD")
	queryString := os.Getenv("QUERY_STRING")

	// Process the request based on method and query
	adminServer.HandleCGIRequest(method, queryString)
}

// startDevServer starts the development server
func startDevServer() {
	// Create public directory if it doesn't exist
	os.MkdirAll("public", 0755)
	os.MkdirAll("public/data", 0755)

	// Serve files from public directory
	http.Handle("/", http.FileServer(http.Dir("./public/")))

	// Handle admin CGI requests
	http.HandleFunc("/admin.cgi", func(w http.ResponseWriter, r *http.Request) {
		// Simulate CGI processing for development
		queryString := r.URL.Query().Encode()
		if r.Method == "POST" {
			// For dev mode, we'll just show a form processing page
			fmt.Fprint(w, "<html><body><h1>Form Submitted</h1><p>This would process your form in CGI mode</p><a href='/admin.cgi'>Back to Admin</a></body></html>")
		} else {
			// Serve the admin home page
			fmt.Fprint(w, "<html><body><h1>Development Admin Interface</h1><p>Query: "+queryString+"</p><p>This simulates the admin interface in development mode.</p><a href='/admin.cgi'>Admin Home</a></body></html>")
		}
	})

	fmt.Printf("Development server starting on :%s\n", devServerPort)
	log.Fatal(http.ListenAndServe(":"+devServerPort, nil))
}

// startFastCGIServer starts the FastCGI server
func startFastCGIServer() {
	// FastCGI implementation
	log.Fatal("FastCGI server not implemented yet")
}