package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/sprengr/Updater/cli"
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

func clientRedirectHome(w http.ResponseWriter) {
	fmt.Fprint(w, "<script>window.location.replace('/')</script>")
}

func main() {
	if cli.VersionFlag() {
		fmt.Print(Version)
		os.Exit(0)
	}

	update.Cleanup()

	done := make(chan bool, 1)

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
		if ok, release := update.Check(Version); ok && update.Install(release) {
			clientRedirectHome(w)
			done <- true
		} else {
			fmt.Fprint(w, "Nothing to install 😟")
		}
	})
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Println("Serving")

	<-done
	log.Println("Shutting down as newer version is running")
}
