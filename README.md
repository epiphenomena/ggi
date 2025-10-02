# GGI (Go Generated Interfaces)
A foundation library for building custom static sites with auto-generated admin interfaces

## Overview

The core idea is that LLMs are pretty good at producing html/css/js with a desired appearance and functionality,
especially when an experienced professional can guide the process. And most websites don't really need more than
a few static pages -- think a small business that just needs: About Us, Portfolio, Contact Us. The challenge is
that having put together this 1-3 static pages, it is intimidating and error prone for a non-technical user to
keep the content up-to-date (requires editing the html), and any changes to common site elements in sync (requires
editing the html in multiple places or maintaining separate files for css and js which then have to be versioned to
break the caching of clients for updates to be visible). So the proposal is that the professional can:

1) work with an LLM to produce the initial site design (with temporary filler content)
2) break the design up into separate css, js and html templates (so that elements common to multiple pages are shared)
3) identify and pull out the handful of elements that the non-technical web admin will want to edit
4) provide for that web admin a very simple, clean admin ui for editing those elements and then publishing the changes
5) Have the final product served to clients be fast "hand written" static html/css/js files

This project will serve as a framework or sample project that the professional clones and then edits to achieve this end.

Beyond html/css/js, the main programming language will be go because go has solid support for core internet protocols,
built in simple templating, compiles to a static binary (for easy deployment), and supports embedding resources in the
binary (so that only a single file needs to be uploaded to the server to update the admin ui and build script).

As an example:

1) the professional could work with LLMs to get a basic html/css/js page together for a small business.
2) The page has sections: hero, about, portfolio (a series of cards with a picture and description), testimonials, and contact us.
3) The professional isolates the different portions of the page into their own template files inside /site. And hooks those templates up in build.go
4) The professional sets the hero image to be replaceable by the user by placing it in the /public/data folder and making the appriopriate changes to the templates and build.go
5) The professional sets the content section to be editable by the user as markdown by creating a placeholder markdown file in /public/data and making the appriopriate changes to the templates and build.go
6) The professional sets the portfolio section up to be editable by the user so that the user can add, delete, rearrange the cards, and edit each card separately where editing a card brings up a form containing a text input for the description and an input to upload a replacement image
7) The professional sets the contact information to be editable as a form with inputs for each part of the contact information, save the data to json
8) The professional compiles the go into admin.cgi
9) uploads the contents of /public to the server
10) Done.

## Project Structure


- /public: Directory to be uploaded to web server
  - /data: Directory, protected from being served with htaccess containing the data files that are edited by the admin ui and used by the build script to produce the public site
    - *.json, *.md, *.jpg, *.png, *.ico, *.mp4, etc.: data files that can be edited / replaced (via upload of new media files) via the admin ui
  - admin.cgi: location of compiled binary. chmod 755. htaccess requires BasicAuth to access. Functionality described further below.
  - .htaccess: contains rules implementing above access restrictions
- /adminui: Directory containing resources for the admin ui: html templates, css, js, images
- /site: Directory containing resources for the public site: html templates, css, js, images
- main.go: entrypoint for cli / cgi script
- build.go: ideally this file contains a simple, clear primary function (calling potentially other functions) that embeds the resources in /site and then reads the files in /public/data to render the public site inside /public. This primary function is called by the cgi binary after any edits to the files in /public/data
- admin.go: embeds the resources in /adminui and pulls together the implementation of the adminui for including in the cgi script
- *.go: other go files it makes sense to break out on their own for clarity

### The compiled admin.cgi supports:

- Calling with cli arg `serve` to start a dev server that:
  - watches all of the relevant files and rebuilds both the admin.cgi and public site on change
  - serves all of the static files in public
  - calls the admin.cgi script as a CGI script with appropriate env variables set and stdin data as needed
- Calling with cli arg `fastcgi` to start a fastcgi compatible server that handles requests just like the cgi
- Calling with no cli args as a cgi script that:
  - responds to GET with no query args -> the home page of the admin ui that provides links for accessing the forms to edit the files in ./data
  - responds to GET with query indicating a ./data/* file -> a form to edit and submit changes to that file
  - responds to POST with query indicating a ./data/* file -> replaces selected file in ./data with form data -> rebuilds site -> reloads GET with query for that same file with a toast to indicate success or failure


### The admin ui:

- every page should have a header with links to the admin home page and to the site
- Should be simple clean and modern, basic, easy to use
- mobile accessible but assume that most edits will take place from a desktop browser
- should be contain everything needed to produce automatically layed out forms for editing the different data files:
  - markdown
  - replace media by uploading
  - edit json by having appropriate inputs for values labelled by their key (and handling nested objects / arrays)

## Getting Started

### Prerequisites
- Go 1.16 or higher (for embed functionality)

### Development
1. Clone the repository
2. Run development server: `go run *.go --serve`
3. Visit http://localhost:8080/admin.cgi to access the admin interface
4. Edit templates in `/site/templates/` and data in `/public/data/` as needed

### Building the Site
- Build the static site: `go run *.go --build`
- The site files will be generated in the `/public` directory

### Building for Production
1. Build the admin CGI binary: `go build -o public/admin.cgi`
2. Make it executable: `chmod +x public/admin.cgi`
3. Upload the contents of `/public` to your web server
4. Configure your web server to handle CGI requests to admin.cgi
5. Set up BasicAuth using the .htaccess file provided

## Implementation Details

### Build Process
The build process:
1. Embeds resources from `/site` directory (CSS, JS, templates)
2. Reads data files from `/public/data/`
3. Processes templates with data to generate static HTML
4. Places generated files in `/public/` directory

### Data File Handling
- JSON files are parsed and made available to templates
- Markdown files are rendered to HTML
- Media files are copied directly to the public directory
- Form submissions update the corresponding data files

### Admin Interface Features
- File browser for data files
- Type-appropriate editors (JSON editor for .json, text editor for .md/.txt)
- Site rebuild after successful updates
- Mobile-responsive design

## Example Workflow

1. Create a new JSON file in `/public/data/` (e.g., `site.json`) with your site data
2. Update `/site/templates/index.tmpl` to use your data variables
3. Run `go run *.go --build` to generate the site
4. The site will be built in the `/public/` directory with your data populated
5. Use the admin interface to edit your data files through the web interface