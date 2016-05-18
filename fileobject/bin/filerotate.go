package main

import (
	"flag"
	"fmt"
	"os"

	"agerotate/bucket"
	"agerotate/fileobject/config"
)

var (
	ConfigPath = flag.String("config", "", "Path to file rotation config.")
	FieldSep   = flag.String("fieldsep", ":", "Field separator for range lines.")
)

func errorExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(-1)
}

func main() {
	flag.Parse()

	config, err := os.Open(*ConfigPath)
	if err != nil {
		errorExit("Error opening config %q: %v\n", *ConfigPath, err)
	}

	files, ranges, err := config.Parse(config, *FieldSep)
	if err != nil {
		errorExit("Error parsing config %q: %v\n", *ConfigPath, err)
	}

	err = bucket.Cleanup(ranges, files)
	if err != nil {
		errorExit("Error doing cleanup: %v\n", err)
	}
}
