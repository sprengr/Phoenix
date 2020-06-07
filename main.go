package main

import (
	"flag"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/sprengr/Updater/update"
)

// const Version = "ver3"
var Version string

var (
	pageTemplate = `
<!DOCTYPE html>
<html>
<head>
<title> Server {{.Version}} </title>
</head>
<body>
<h1>This server is version {{.Version}}</h1>
<a href="check">Check for new version</a>
<br>
{{if .NewVersion}}New version is available: {{.NewVersion}} | <a
href="install">Upgrade</a>{{end}}
</body>
</html>
`
	Status = struct{ Version, NewVersion string }{Version, ""}
)

func checkVersionFlag() {
	var flgVersion bool
	flag.BoolVar(&flgVersion, "version", false, "if true, print version and exit")
	flag.Parse()
	if flgVersion {
		fmt.Printf(Version)
		os.Exit(0)
	}
}

func main() {
	checkVersionFlag()
	page, err := template.New("page").Parse(pageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := page.Execute(w, Status); err != nil {
			log.Fatal(err)
		}
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if ok, release := update.Check(Version); ok {
			fmt.Fprintf(w, "New version is available: %v", release.Version)
			Status.NewVersion = release.Version
		} else {
			fmt.Fprint(w, "You are already on the latest version 🚀")
		}
	})
	http.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Not implemented %v", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
