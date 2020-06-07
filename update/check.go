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
func asOld(process string) string {
	return process + ".old"
}

func getVersion(process string) (bool, string) {
	cmd := exec.Command(asUpdate(process), "--version")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return true, string(out)
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
	return true
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
		log.Printf("Previous version found (%v) but we ignore it as it's not a regular file")
		return
	}

	err = os.Remove(asOld(currentVersion))

	if err != nil {
		log.Fatal("Couldn't remove previous version: %v", err)
		return
	}
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
