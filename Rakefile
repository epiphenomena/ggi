require 'rake'

# Default task
task :default => :build

desc "Upload to production"
task :upload do
  sh "rsync -avzh public/ dest/"
end

# Build the admin CGI binary
desc "Build the admin CGI binary and site"
task :build => [:cgi, :site]

desc "Build the admin CGI binary"
task :cgi do
  sh "go build -o public/admin.cgi ."
  sh "chmod +x public/admin.cgi"
end

# Clean the public folder
desc "Clean the public folder of build artifacts"
task :clean => :cgi do
  sh "public/admin.cgi --clean"
end

# Build and run the development server
desc "Run the development server"
task :serve => :build do
  sh "public/admin.cgi --serve"
end

# Build and run a test of the CGI script
desc "Test the CGI script"
task :cgitest => :build do
  puts "Testing CGI script..."
  ENV['REQUEST_METHOD'] = 'GET'
  ENV['QUERY_STRING'] = ''
  sh "./public/admin.cgi"
end

# Run the build command
desc "Build the static site"
task :site => :cgi do
  sh "public/admin.cgi --build"
end

# Run all commands in sequence
desc "Clean, build, and serve"
task :dev => [:clean, :build, :serve]