package ggi

import (
	"html/template"
	"net/http"
	"os"
)

// SecretKey holds the authentication token
var SecretKey string

// CheckSecretKey checks if the secret key has been set, and panics if not
func checkSecretKey() {
	if SecretKey == "" {
		panic("SecretKey must be set. Please set ggi.SecretKey to a secure value before using this package.")
	}
}

// IsCGI checks if the program is running as a CGI script
func IsCGI() bool {
	// Check for common CGI environment variables
	_, hasRequestMethod := os.LookupEnv("REQUEST_METHOD")
	_, hasScriptName := os.LookupEnv("SCRIPT_NAME")
	
	return hasRequestMethod && hasScriptName
}

// IsAuthenticated checks if the request contains the correct secret key
func IsAuthenticated(r *http.Request) bool {
	// Check that secret key has been set first
	checkSecretKey()
	
	if r.Method != "POST" {
		return false
	}
	
	token := r.FormValue("token")
	return token == SecretKey
}



// BaseTemplate defines the base HTML template
const BaseTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    {{block "css" .}}{{end}}
    {{block "js" .}}{{end}}
</head>
<body>
    <div class="container">
        {{block "header" .}}
        <header>
            <h1>{{.Title}}</h1>
        </header>
        {{end}}
        
        {{block "content" .}}
        <main>
            {{.Content}}
        </main>
        {{end}}
        
        {{block "footer" .}}
        <footer>
            <p>&copy; 2025 GGI Site</p>
        </footer>
        {{end}}
    </div>

    <!-- Secret Key Modal -->
    <div id="secretKeyModal" class="modal">
        <div class="modal-content">
            <span class="close">&times;</span>
            <h2>Settings</h2>
            <p>Enter your secret key to manage content:</p>
            
            <div class="secret-key-form">
                <label for="secretKey">Secret Key:</label><br>
                <input type="password" id="secretKey" name="secretKey" value="">
                <br>
                <button onclick="saveSecretKey()">Save Key</button>
                <button onclick="clearSecretKey()">Clear Key</button>
                
                <div style="margin-top: 20px;">
                    <h3>About GGI</h3>
                    <p>GGI is a Go-based library that combines a simple static site generator with a CGI-based editor.</p>
                    <p>Use the secret key to authenticate requests for editing content.</p>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Get modal element
        var modal = document.getElementById("secretKeyModal");
        
        // Get span element that closes the modal
        var span = document.getElementsByClassName("close")[0];
        
        // When the user clicks on the span (x), close the modal
        span.onclick = function() {
            modal.style.display = "none";
        }
        
        // When the user clicks anywhere outside of the modal, close it
        window.onclick = function(event) {
            if (event.target == modal) {
                modal.style.display = "none";
            }
        }
        
        // Load saved secret key from localStorage
        document.addEventListener("DOMContentLoaded", function() {
            var savedKey = localStorage.getItem("ggi_secret_key");
            if (savedKey) {
                document.getElementById("secretKey").value = savedKey;
            }
        });
        
        // Save secret key to localStorage
        function saveSecretKey() {
            var key = document.getElementById("secretKey").value;
            if (key) {
                localStorage.setItem("ggi_secret_key", key);
                alert("Secret key saved to localStorage!");
                modal.style.display = "none";
            } else {
                alert("Please enter a secret key.");
            }
        }
        
        // Clear secret key from localStorage
        function clearSecretKey() {
            localStorage.removeItem("ggi_secret_key");
            document.getElementById("secretKey").value = "";
            alert("Secret key cleared!");
        }
        
        // Add secret key to forms automatically
        document.addEventListener("DOMContentLoaded", function() {
            // Add secret key to all forms
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
</body>
</html>
`

// ParseBaseTemplate parses the base template
func ParseBaseTemplate() (*template.Template, error) {
	return template.New("base").Parse(BaseTemplate)
}

// ProtectedHandler wraps handlers that require authentication
func ProtectedHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Error(w, "Unauthorized: Invalid or missing token", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

// ParseTemplate parses a template with the base layout
func ParseTemplate(tmpl string) (*template.Template, error) {
	baseTmpl, err := ParseBaseTemplate()
	if err != nil {
		return nil, err
	}
	
	// Create a new template based on the base
	t := template.Must(baseTmpl.Clone())
	return template.Must(t.New("page").Parse(tmpl)), nil
}

// ParseTemplateWithCSS parses a template with CSS configuration
func ParseTemplateWithCSS(tmpl string, cssConfig interface{}) (*template.Template, error) {
	// Parse the base template
	baseTmpl, err := ParseBaseTemplate()
	if err != nil {
		return nil, err
	}
	
	// Combine the templates
	combinedTmpl := template.Must(baseTmpl.New("css").Parse(CSSTemplate))
	t := template.Must(combinedTmpl.Clone())
	
	// Add the page template
	return template.Must(t.New("page").Parse(tmpl)), nil
}