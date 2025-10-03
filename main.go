package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
		// Build the admin.cgi binary first if it doesn't exist
		if _, err := os.Stat("public/admin.cgi"); os.IsNotExist(err) {
			cmd := exec.Command("go", "build", "-o", "public/admin.cgi", ".")
			if err := cmd.Run(); err != nil {
				http.Error(w, "Failed to build admin.cgi: "+err.Error(), 500)
				return
			}
		}

		// Set up environment variables for CGI
		env := append(os.Environ(),
			"REQUEST_METHOD="+r.Method,
			"QUERY_STRING="+r.URL.RawQuery,
			"SCRIPT_NAME=/admin.cgi",
			"SERVER_PROTOCOL="+r.Proto,
			"HTTP_HOST="+r.Host,
			"SERVER_SOFTWARE=GGI-Dev-Server",
			"SERVER_NAME=localhost",
			"SERVER_PORT="+devServerPort,
		)

		// Create the command
		cmd := exec.Command("./public/admin.cgi")
		cmd.Env = env
		cmd.Dir = "public"

		// Capture the output
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, "CGI script error: "+err.Error()+"\nOutput: "+string(output), 500)
			return
		}

		// Write the output to the HTTP response
		// Parse the headers from the CGI output
		responseStr := string(output)
		
		// Find the end of headers (double newline)
		parts := strings.SplitN(responseStr, "\r\n\r\n", 2)
		if len(parts) == 2 {
			// Parse headers
			headers := strings.Split(parts[0], "\r\n")
			for _, header := range headers {
				if colonIndex := strings.Index(header, ":"); colonIndex > 0 {
					headerName := strings.TrimSpace(header[:colonIndex])
					headerValue := strings.TrimSpace(header[colonIndex+1:])
					w.Header().Set(headerName, headerValue)
				}
			}
			// Write the body
			fmt.Fprint(w, parts[1])
		} else {
			// If no headers, just write the content with a default content type
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, responseStr)
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