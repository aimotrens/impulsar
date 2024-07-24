package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/aimotrens/impulsar/cout"
)

var (
	compileDate     string
	impulsarVersion string
)

func version(fl flagLoader) {
	fl(func(fs *flag.FlagSet) {
		fs.Usage = noFlagsUsage("version")
	})

	v := impulsarVersion
	if v == "" {
		v = "x.x.x"
	}

	cd := compileDate
	if cd == "" {
		cd = "unknown"
	}

	fmt.Println(cout.Bold("Impulsar " + v))
	fmt.Printf("Compiled at %s with %s\n\n", cd, cout.Cyan(runtime.Version()))
	os.Exit(0)
}
