package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sprengr/Phoenix/cli"
	"github.com/sprengr/Phoenix/render"
	"github.com/sprengr/Phoenix/update"
)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render.Index(w, Status)
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if release, ok := update.Check(Version); ok {
			Status.VersionFound = release.Version
		}
		render.Check(w, Status)
	})

	http.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		release, ok := update.Check(Version)
		if ok && update.Install(release) {
			Status.VersionInstalled = release.Version
		}
		render.Install(w, Status)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Println("Serving")

	checkPeriodically()

	<-update.Shutdown
	log.Println("Shutting down as newer version is running")
}

func checkPeriodically() {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if release, ok := update.Check(Version); ok {
					Status.VersionFound = release.Version
				}
			}
		}
	}()
}
