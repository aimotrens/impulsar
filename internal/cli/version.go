package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/aimotrens/impulsar/pkg/tui"
)

func version(fl flagLoader, buildInfo BuildInfoProvider) {
	fl(func(fs *flag.FlagSet) {
		fs.Usage = noFlagsUsage("version")
	})

	compileDate, impulsarVersion := buildInfo()
	compileTime := "unknown"

	if seconds, err := strconv.ParseInt(compileDate, 10, 64); err == nil {
		compileTime = time.Unix(seconds, 0).Format(time.RFC1123)
	}

	fmt.Println(tui.Bold("Impulsar " + impulsarVersion))
	fmt.Printf("Compiled at %s with %s\n\n", compileTime, tui.Cyan(runtime.Version()))
	os.Exit(0)
}
