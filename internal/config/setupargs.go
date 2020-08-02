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
	flag.Parse()

	if showVersion == true {
		fmt.Printf("Package version: v%[1]s\nYAML Config Version: v%[2]d\n", AppConfig.BuildVersion, AppConfig.ConfigVersion)
		os.Exit(1)
	}

	if YamlPath == "" {
		fmt.Printf("Config file not provided. See '%[1]s --help'\n", os.Args[0])
		os.Exit(1)
	}
}
