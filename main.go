package main

import (
	"flag"
	"fmt"
	"github.com/aimotrens/impulsar/cout"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/engine"
	_ "github.com/aimotrens/impulsar/engine/docker_executor"
	_ "github.com/aimotrens/impulsar/engine/shell_executor"
	_ "github.com/aimotrens/impulsar/engine/ssh_executor"

	"github.com/aimotrens/impulsar/model"
	"gopkg.in/yaml.v3"

	_ "embed"
)

var (
	//go:embed shell-completion/bash.sh
	bashCompletion string

	//go:embed shell-completion/zsh.sh
	zshCompletion string

	//go:embed shell-completion/powershell.ps1
	powershellCompletion string

	compileDate     string
	impulsarVersion string
)

func main() {
	var impulsarFile string
	var envVars model.FlagArray
	var dumpJobs bool
	var showJobs bool
	var genBashCompletion bool
	var genZshCompletion bool
	var genPowershellCompletion bool
	var version bool

	flag.BoolVar(&version, "v", false, "version")
	flag.BoolVar(&version, "version", false, "version")

	flag.StringVar(&impulsarFile, "f", "./impulsar.yml", "impulsar file")
	flag.Var(&envVars, "e", "additional environment variables")
	flag.BoolVar(&dumpJobs, "dump-jobs", false, "dump parsed jobs to impulsar-dump.yml")
	flag.BoolVar(&showJobs, "show-jobs", false, "show jobs from impulsar.yml")

	flag.BoolVar(&genBashCompletion, "gen-bash-completion", false, "generate bash completion")
	flag.BoolVar(&genZshCompletion, "gen-zsh-completion", false, "generate zsh completion")
	flag.BoolVar(&genPowershellCompletion, "gen-powershell-completion", false, "generate powershell completion")

	flag.Parse()

	if version {
		fmt.Printf("Impulsar %s\nCompiled at %s with %s\n\n", impulsarVersion, compileDate, runtime.Version())
		os.Exit(0)
	}

	if genBashCompletion {
		fmt.Println(bashCompletion)
		os.Exit(0)
	}

	if genZshCompletion {
		fmt.Println(zshCompletion)
		os.Exit(0)
	}

	if genPowershellCompletion {
		fmt.Println(powershellCompletion)
		os.Exit(0)
	}

	if showJobs {
		impulsar := loadimpulsarFile(impulsarFile)
		for key := range impulsar {
			fmt.Println(key)
		}
		os.Exit(0)
	}

	fmt.Printf(cout.Bold(fmt.Sprintf("Impulsar %s\n", impulsarVersion)))

	if flag.NArg() == 0 {
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

	impulsar := loadimpulsarFile(impulsarFile)
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

	fmt.Println(cout.DarkGray("Execution plan ..."))
	for i := 0; i < flag.NArg(); i++ {
		fmt.Println(cout.DarkGray("- " + flag.Arg(i)))
	}
	fmt.Println("")

	// Alle Job-Argumente sammeln
	for i := 0; i < flag.NArg(); i++ {
		e.CollectArgs(flag.Arg(i))
	}

	// Alle Jobs ausführen
	for i := 0; i < flag.NArg(); i++ {
		e.RunJob(flag.Arg(i))
	}
}

func loadimpulsarFile(impulsarFile string) map[string]*model.Job {
	f, err := os.Open(impulsarFile)
	if err != nil {
		panic(err)
	}

	yData, _ := io.ReadAll(f)

	var impulsar map[string]*model.Job
	dec := yaml.NewDecoder(strings.NewReader(string(yData)))
	dec.KnownFields(true)

	err = dec.Decode(&impulsar)
	if err != nil {
		panic(err)
	}

	return impulsar
}
