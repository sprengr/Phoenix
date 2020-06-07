package cli

import (
	"flag"
)

func VersionFlag() bool {
	var flgVersion bool
	flag.BoolVar(&flgVersion, "version", false, "if true, print version and exit")
	flag.Parse()

	// if flgVersion {
	// 	fmt.Print(Version)
	// }

	return flgVersion
}
