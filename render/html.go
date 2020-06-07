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

func Index(w http.ResponseWriter, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, "index.gohtml", data); err != nil {
		log.Fatal(err)
	}
}

func Check(w http.ResponseWriter, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, "check.gohtml", data); err != nil {
		log.Fatal(err)
	}
}

func Install(w http.ResponseWriter, data interface{}) {
	if err := tmpl.ExecuteTemplate(w, "installSucess.gohtml", data); err != nil {
		log.Fatal(err)
	}
}
