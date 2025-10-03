require 'rake'

# Default task
task :default => :build

# Build the admin CGI binary
desc "Build the admin CGI binary"
task :build do
  system("go build -o public/admin.cgi .") || abort("Build failed")
  system("chmod +x public/admin.cgi") || abort("Failed to make executable")
  puts "Admin CGI built successfully!"
end

# Clean the public folder
desc "Clean the public folder of build artifacts"
task :clean do
  system("go run *.go --clean") || abort("Clean failed")
end

# Build and run the development server
desc "Run the development server"
task :serve => :build do
  system("go run *.go --serve") || abort("Serve failed")
end

# Build and run a test of the CGI script
desc "Test the CGI script"
task :cgitest => :build do
  puts "Testing CGI script..."
  system("REQUEST_METHOD='GET' QUERY_STRING='' ./public/admin.cgi")
end

# Run the build command
desc "Build the static site"
task :site => :build do
  system("go run *.go --build") || abort("Site build failed")
end

# Run all commands in sequence
desc "Clean, build, and serve"
task :dev => [:clean, :build, :serve]