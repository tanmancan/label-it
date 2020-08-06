package config

import (
	"flag"
	"fmt"
	"os"
)

// YamlPath path to the yaml config file provided via a flag
var YamlPath string

// ShowHelp value of flag to show help
var ShowHelp bool

// DryRun outputs labels to be added to pr based on rules without making an API call
var DryRun bool

// SetupArgs sets up flags and help text
func SetupArgs() {
	flag.Usage = func() {
		fmt.Printf("Usage: %[1]s [--version][--help][-c <path>]\n", os.Args[0])
		fmt.Printf("Example: %[1]s -c label-it.yaml\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	var showVersion bool

	flag.BoolVar(&ShowHelp, "-help", false, "Display the help text")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.StringVar(&YamlPath, "c", "", "Path to the yaml file")
	flag.BoolVar(&DryRun, "dry", false, "Outputs list of pull request and matched labels. Does not call the API")
	flag.Parse()

	if showVersion == true {
		fmt.Printf("Package version: v%[1]s\nYAML Config Version: v%[2]d\n", BuildVersion, ConfigVersion)
		os.Exit(1)
	}

	if YamlPath == "" {
		fmt.Printf("Config file not provided. See '%[1]s --help'\n", os.Args[0])
		os.Exit(1)
	}
}
