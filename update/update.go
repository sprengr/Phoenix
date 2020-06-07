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

var Shutdown chan bool = make(chan bool, 1)

type Release struct {
	Version    string
	executable string
}

func Cleanup() {
	currentVersion, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return
	}

	oldVersionStat, err := os.Stat(asOld(currentVersion))
	if os.IsNotExist(err) {
		return
	}

	if !oldVersionStat.Mode().IsRegular() {
		log.Printf("Previous version found but we ignore it as it's not a regular file")
		return
	}

	err = os.Remove(asOld(currentVersion))

	if err != nil {
		log.Fatalf("Couldn't remove previous version: %v", err)
		return
	}
}

func Check(currentVersion string) (Release, bool) {
	executable, ok := getExecutable()

	if !ok {
		return Release{}, false
	}

	if version, ok := getVersion(executable); ok && version != currentVersion {
		return Release{version, executable}, true
	}

	return Release{}, false
}

func getExecutable() (string, bool) {
	executables, err := ioutil.ReadDir(releasesPath)
	if err != nil {
		log.Fatal(err)
		return "", false
	}

	currentExecutable, err := os.Executable()

	if err != nil {
		log.Fatal(err)
		return "", false
	}

	for _, e := range executables {
		executable := e.Name()
		if executable == filepath.Base(currentExecutable) {
			log.Printf("Found new executable %v", executable)
			return executable, true
		}
	}
	return "", false
}

func getVersion(process string) (string, bool) {
	cmd := exec.Command(asUpdate(process), "--version")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
		return "", false
	}

	return string(out), true
}

func asUpdate(process string) string {
	return releasesPath + "/" + process
}

func Install(release Release) bool {
	currentVersion, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return false
	}

	err = os.Rename(currentVersion, asOld(currentVersion))
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Printf("Renamed current executable to %v", asOld(currentVersion))

	nBytes, err := copy(asUpdate(release.executable), currentVersion)
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Printf("Copied new version %v (%d)", currentVersion, nBytes)

	err = start(currentVersion)
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Printf("Successfully updated")
	Shutdown <- true
	return true
}

func asOld(process string) string {
	return process + ".old"
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func start(process string) error {
	_, err := os.Stat(process)
	if os.IsNotExist(err) {
		return err
	}

	cmd := exec.Command(process)
	return cmd.Start()
}
