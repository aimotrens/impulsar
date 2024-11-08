package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aimotrens/impulsar/cout"

	"github.com/aimotrens/impulsar/engine"
	_ "github.com/aimotrens/impulsar/engine/docker_executor"
	_ "github.com/aimotrens/impulsar/engine/shell_executor"
	_ "github.com/aimotrens/impulsar/engine/ssh_executor"

	"github.com/aimotrens/impulsar/model"
	"gopkg.in/yaml.v3"
)

func run(fl flagLoader, buildInfo BuildInfoProvider) {
	_, impulsarVersion := buildInfo()

	var impulsarFile string
	var envVars model.FlagArray
	var dumpJobs bool
	var dryrun bool

	runFlags := fl(func(fs *flag.FlagSet) {
		fs.Usage = func() {
			fmt.Println("Usage: impulsar " + ColorizeCmd("run") + " [-f <impulsar file>] [-e <name=value>]... [-dump-jobs] <job>...")
			fs.PrintDefaults()
		}

		fs.StringVar(&impulsarFile, "f", "./impulsar.yml", "impulsar file")
		fs.Var(&envVars, "e", "additional environment variables (format: name=value)")
		fs.BoolVar(&dumpJobs, "dump-jobs", false, "dump parsed jobs to impulsar-dump.yml")
		fs.BoolVar(&dryrun, "dryrun", false, "dryrun, only show execution plan")
	})

	fmt.Print(cout.Bold(fmt.Sprintf("Impulsar %s\n", impulsarVersion)))

	if runFlags.NArg() == 0 {
		fmt.Println(cout.Yellow("No jobs provided"))
		return
	}

	addtitionalEnvVars := make(model.VariableMap)
	for _, v := range envVars {
		if !strings.Contains(v, "=") {
			fmt.Println(cout.Red("Invalid variable: " + v))
			os.Exit(1)
		}

		kv := strings.Split(v, "=")
		os.Setenv(kv[0], kv[1])
		addtitionalEnvVars[kv[0]] = kv[1]
	}

	impulsar := loadImpulsarFile(impulsarFile)
	for key, job := range impulsar {
		job.Name = key
		job.SetDefaults()
	}

	if dumpJobs {
		dump, _ := yaml.Marshal(impulsar)
		f, _ := os.Create("./impulsar-dump.yml")
		defer f.Close()
		f.Write(dump)
	}

	e := engine.New(impulsar, addtitionalEnvVars)

	var plan string
	var canExecute bool
	fmt.Print(cout.FormattingArea(nil, func(b *strings.Builder) {
		fmt.Fprintln(b, "Execution plan ...")

		plan, canExecute = buildExecutionPlan(e, runFlags.Args())
		fmt.Fprintln(b, plan)
		fmt.Fprintln(b)

		fmt.Fprintln(b, "")
	}))

	if !canExecute {
		fmt.Println(cout.Red("Execution plan is invalid, see above"))
		os.Exit(1)
	}

	if dryrun {
		return
	}

	// collect all arguments
	for i := 0; i < runFlags.NArg(); i++ {
		e.CollectArgs(runFlags.Arg(i))
	}

	// run all jobs
	for i := 0; i < runFlags.NArg(); i++ {
		e.RunJob(runFlags.Arg(i))
	}
}

func loadImpulsarFile(impulsarFile string) map[string]*model.Job {
	var f *os.File
	var err error

	if impulsarFile == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(impulsarFile)
		if err != nil {
			panic(err)
		}
	}

	var impulsar map[string]*model.Job
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)

	err = dec.Decode(&impulsar)
	if err != nil {
		panic(err)
	}

	return impulsar
}
