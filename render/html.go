package render

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("render/template/*.gohtml"))
}

// Index renders the initial page into w.
func Index(w http.ResponseWriter, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, "index.gohtml", data); err != nil {
		log.Fatal(err)
	}
}

// Check renders the page informing if there's an update available into w.
func Check(w http.ResponseWriter, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, "check.gohtml", data); err != nil {
		log.Fatal(err)
	}
}

// Install renders the installation result into w.
func Install(w http.ResponseWriter, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, "install.gohtml", data); err != nil {
		log.Fatal(err)
	}
}
