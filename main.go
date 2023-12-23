package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/engine"
	"github.com/aimotrens/impulsar/model"
	"gopkg.in/yaml.v3"

	_ "embed"
)

var (
	//go:embed bash-autocompletion.sh
	bashAutoComplete string
	compileDate      string
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var impulsarFile string
	var envVars arrayFlags
	var dumpJobs bool
	var showJobs bool
	var genBashCompletion bool

	flag.StringVar(&impulsarFile, "f", "./impulsar.yml", "impulsar file")
	flag.Var(&envVars, "e", "environment variables")
	flag.BoolVar(&dumpJobs, "dump-jobs", false, "verbose")
	flag.BoolVar(&showJobs, "show-jobs", false, "verbose")
	flag.BoolVar(&genBashCompletion, "gen-bash-completion", false, "verbose")

	flag.Parse()

	if genBashCompletion {
		fmt.Println(bashAutoComplete)
		os.Exit(0)
	}

	fmt.Printf("Impulsar\nCompiled at %s with %s\n\n", compileDate, runtime.Version())

	if showJobs {
		impulsar := loadimpulsarFile(impulsarFile)
		for key := range impulsar {
			fmt.Println(key)
		}
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("No jobs provided")
		return
	}

	envVarsMap := make(map[string]string)
	for _, v := range envVars {
		if !strings.Contains(v, "=") {
			fmt.Println("Invalid variable:", v)
		}

		kv := strings.Split(v, "=")
		os.Setenv(kv[0], kv[1])
		envVarsMap[kv[0]] = kv[1]
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

	e := engine.New(impulsar, envVarsMap)

	fmt.Println("Execution plan ...")
	for i := 0; i < flag.NArg(); i++ {
		fmt.Println("-", flag.Arg(i))
	}
	fmt.Println("")

	for i := 0; i < flag.NArg(); i++ {
		e.RunJob(flag.Arg(i))
	}
}

func loadimpulsarFile(impulsarFile string) model.ImpulsarList {
	f, err := os.Open(impulsarFile)
	if err != nil {
		panic(err)
	}

	yData, _ := io.ReadAll(f)

	var impulsar model.ImpulsarList
	dec := yaml.NewDecoder(strings.NewReader(string(yData)))
	dec.KnownFields(true)

	err = dec.Decode(&impulsar)
	if err != nil {
		panic(err)
	}

	return impulsar
}
