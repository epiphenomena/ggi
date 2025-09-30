# ggi
Go lang library that combines a simple static site generator with a cgi based editor

The fastest and easiest site to serve is static html. But it's complicated for normal users to edit and does not support templating elements that are repeated within a page or across pages. 

SSGs are a solution to the templating problem, but require learning and conforming to their idiosyncrasies. 

There are admin UI's (notably Wordpress) that attempt to make it easier for users to edit content without needing to understand the tech stack.

However, both SSGs adn the admin UIs necessarily develop a great deal of complexity in order to support a wide range of use cases.

LLMs make creating customized websites and customized admin UIs easy. The goal of this project is to create a simple go lang based library to support an LLM driven website creation and maintenance.

The idea is to import this library into a new website project, add the needed customizations for that project, and then compile to single binary that acts as a cgi script.

The CGI script supports editing source files and then generating the resulting static html.
