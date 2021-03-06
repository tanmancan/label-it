package config

import (
	"errors"
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

// AutoConfirm value of flag to used to auto confirm any prompt
var AutoConfirm bool

// SetupArgs sets up flags and help text
func SetupArgs() error {
	flag.Usage = func() {
		fmt.Printf("Usage: %[1]s [--version][--help][-c <path>]\n", os.Args[0])
		fmt.Printf("Example: %[1]s -c label-it.yaml\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	var showVersion bool

	flag.BoolVar(&ShowHelp, "-help", false, "Display the help text")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.StringVar(&YamlPath, "c", "", "Path to the yaml file")
	flag.BoolVar(&DryRun, "dry", false, "Outputs list of pull request and matched labels. Does not call the API")
	flag.BoolVar(&AutoConfirm, "y", false, "Auto confirms user prompt")
	flag.Parse()

	if showVersion == true {
		fmt.Printf("Version: %[1]s\nAPI Version: %[2]s\nSHA: %[3]s\n", BuildVersion, APIVersion, GitSHA)
		os.Exit(0)
	}

	if YamlPath == "" {
		errMessage := fmt.Sprintf("Config file not provided. See '%[1]s --help'\n", os.Args[0])
		return errors.New(errMessage)
	}

	return nil
}
