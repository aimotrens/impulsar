package cli

import (
	"flag"
	"fmt"
)

func showJobs(fl flagLoader) {
	var impulsarFile string

	fl(func(fs *flag.FlagSet) {
		fs.Usage = func() {
			fmt.Println("Usage: impulsar " + ColorizeCmd("show-jobs") + " [-f <impulsar file>]")
			fs.PrintDefaults()
		}
		fs.StringVar(&impulsarFile, "f", "./impulsar.yml", "impulsar file")
	})

	impulsar := loadImpulsarFile(impulsarFile)
	for key := range impulsar {
		fmt.Println(key)
	}
}
