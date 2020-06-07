package update

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const releasesPath = "./releases"

type Release struct {
	Version    string
	executable string
}

func Check(currentVersion string) (bool, Release) {
	ok, release := checkRelease()

	if !ok {
		return false, Release{}
	}

	if ok, version := getVersion(release); ok && version != currentVersion {
		return true, Release{version, release}
	}

	return false, Release{}
}

func checkRelease() (bool, string) {
	releases, err := ioutil.ReadDir(releasesPath)
	if err != nil {
		log.Fatal(err)
		return false, ""
	}

	currentVersion, err := os.Executable()

	if err != nil {
		log.Fatal(err)
		return false, ""
	}

	for _, r := range releases {
		release := r.Name()
		if release == filepath.Base(currentVersion) {
			log.Printf("Found release %v", release)
			return true, release
		}
	}
	return false, ""
}

func asUpdate(process string) string {
	return releasesPath + "/" + process
}

func getVersion(process string) (bool, string) {
	cmd := exec.Command(asUpdate(process), "--version")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return true, string(out)
}

