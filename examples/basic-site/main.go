package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ggi/pkg/ggi"
)

// Set secret key before any other initialization occurs
func init() {
	// For development, we can set the secret key from environment
	if secret := os.Getenv("GGI_SECRET_KEY"); secret != "" {
		ggi.SecretKey = secret
	} else {
		// Use a default for development
		ggi.SecretKey = "dev_secret_key_for_testing"
	}
}

// Example struct for form data
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func main() {
	if ggi.IsCGI() {
		// Handle as CGI script
		handleCGI()
	} else {
		// Start development server
		handleDevServer()
	}
}

func handleCGI() {
	// Only handle POST requests
	requestMethod := os.Getenv("REQUEST_METHOD")
	if requestMethod != "POST" {
		fmt.Println("Status: 405 Method Not Allowed")
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("Only POST requests are allowed")
		return
	}

	// Read and parse the request
	body := os.Getenv("CONTENT_LENGTH")
	if body == "" {
		body = "0"
	}
	// For simplicity in this example, we'll just return a success message
	fmt.Println("Content-Type: text/html")
	fmt.Println()
	fmt.Println("<html><body><h1>CGI Request Processed</h1><p>Token authenticated successfully.</p></body></html>")
}

func handleDevServer() {
	// For development, we'll set the secret key automatically
	ggi.SecretKey = "dev_secret_key_for_testing"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/edit-contact" {
			if r.Method == "POST" {
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Error parsing form", http.StatusBadRequest)
					return
				}

				token := r.FormValue("token")
				if token != ggi.SecretKey {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}

				// Process form data
				name := r.FormValue("Name")
				email := r.FormValue("Email")
				phone := r.FormValue("Phone")

				// In a real implementation, save to JSON file
				contact := Contact{Name: name, Email: email, Phone: phone}
				err = ggi.SaveData([]Contact{contact}, "examples/basic-site/_source/data/contacts.json")
				if err != nil {
					http.Error(w, "Error saving data", http.StatusInternalServerError)
					return
				}

				// Render success page
				ggi.RenderMarkdownPage(w, "Contact Saved", fmt.Sprintf("# Contact Saved\n\nName: %s\n\nEmail: %s\n\nPhone: %s", name, email, phone), nil)
			} else {
				// Show form
				contact := Contact{}
				ggi.RenderDataFormPage(w, "Edit Contact", "/edit-contact", contact)
			}
		} else if r.URL.Path == "/edit-markdown" {
			if r.Method == "POST" {
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Error parsing form", http.StatusBadRequest)
					return
				}

				token := r.FormValue("token")
				if token != ggi.SecretKey {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}

				// Save markdown content
				content := r.FormValue("content")
				err = ggi.SaveMarkdown(content, "examples/basic-site/_source/markdown/home.md")
				if err != nil {
					http.Error(w, "Error saving markdown", http.StatusInternalServerError)
					return
				}

				// Render the saved content
				ggi.RenderMarkdownPage(w, "Markdown Saved", content, nil)
			} else {
				// Load existing content
				content, err := ggi.LoadMarkdown("examples/basic-site/_source/markdown/home.md")
				if err != nil {
					content = "# Welcome\n\nEdit this content..."
				}

				// Show markdown editor
				showMarkdownEditor(w, content)
			}
		} else if r.URL.Path == "/settings" {
			// This page will show the settings modal automatically
			showSettingsPage(w)
		} else {
			// Serve a simple test page
			// The base template already includes the modal, so we just need to trigger it
			homeContent := "# Welcome to GGI Example Site\n\nThis is a test page.\n\n- [Edit a contact form](/edit-contact)\n- [Edit markdown content](/edit-markdown)\n- [Settings](/settings) (click to open modal)"
			ggi.RenderMarkdownPage(w, "GGI Example Site", homeContent, nil)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Development server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func showMarkdownEditor(w http.ResponseWriter, content string) {
	// The base template already includes the modal, so we just need to include it in our template
	tmpl := `
{{define "content"}}
<h2>Edit Markdown Content</h2>
<form method="post" action="/edit-markdown">
    <div class="form-field">
        <textarea name="content" rows="20" cols="80">` + content + `</textarea>
    </div>
    <input type="hidden" name="token" value="">
    <input type="submit" value="Save Markdown">
</form>
{{end}}

{{define "js"}}
<script>
// Re-include the modal functionality if needed
document.addEventListener("DOMContentLoaded", function() {
    // Add secret key to forms automatically
    var forms = document.getElementsByTagName("form");
    for (var i = 0; i < forms.length; i++) {
        var form = forms[i];
        if (!form.querySelector("[name='token']")) {
            var tokenInput = document.createElement("input");
            tokenInput.type = "hidden";
            tokenInput.name = "token";
            tokenInput.value = localStorage.getItem("ggi_secret_key") || "";
            form.appendChild(tokenInput);
        }
    }

    // Add event listener to all submit buttons to ensure token is current
    var submitButtons = document.querySelectorAll("input[type='submit'], button[type='submit']");
    for (var i = 0; i < submitButtons.length; i++) {
        submitButtons[i].addEventListener("click", function(e) {
            var form = e.target.closest("form");
            if (form) {
                var tokenInput = form.querySelector("[name='token']");
                if (tokenInput) {
                    tokenInput.value = localStorage.getItem("ggi_secret_key") || "";
                }
            }
        });
    }
});
</script>
{{end}}
`

	parsedTmpl, err := ggi.ParseTemplate(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	context := map[string]interface{}{
		"Title": "Edit Markdown",
	}

	err = parsedTmpl.ExecuteTemplate(w, "base", context)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func showSettingsPage(w http.ResponseWriter) {
	// The modal is already part of the base template, so we just need to trigger it
	tmpl := `
{{define "content"}}
<h2>Settings</h2>
<p>Click the gear icon or link below to open the settings modal to manage your secret key:</p>
<p><a href="javascript:void(0)" onclick="document.getElementById('secretKeyModal').style.display='block'">Open Settings Modal</a></p>
{{end}}

{{define "js"}}
<script>
// Open the modal when the page loads
document.addEventListener("DOMContentLoaded", function() {
    document.getElementById('secretKeyModal').style.display='block';
});
</script>
{{end}}
`

	parsedTmpl, err := ggi.ParseTemplate(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	context := map[string]interface{}{
		"Title": "Settings",
	}

	err = parsedTmpl.ExecuteTemplate(w, "base", context)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}