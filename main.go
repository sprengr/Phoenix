package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sprengr/Updater/cli"
	"github.com/sprengr/Updater/render"
	"github.com/sprengr/Updater/update"
)

// const Version = "ver3"

var (
	Version string
	Status  = struct{ Version, VersionFound, VersionInstalled string }{Version, "", ""}
)

func main() {
	if cli.VersionFlag() {
		fmt.Print(Version)
		os.Exit(0)
	}

	update.Cleanup()

	done := make(chan bool, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render.Index(w, Status)
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if ok, release := update.Check(Version); ok {
			Status.VersionFound = release.Version
		}
		render.Check(w, Status)
	})
	http.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		if ok, release := update.Check(Version); ok && update.Install(release) {
			done <- true
			Status.VersionInstalled = release.Version
		}
		render.Install(w, Status)
	})
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Println("Serving")
	<-done
	log.Println("Shutting down as newer version is running")
}
