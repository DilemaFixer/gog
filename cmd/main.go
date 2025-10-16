package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DilemaFixer/gog/internal/api"
	"github.com/DilemaFixer/gog/internal/log"
	"github.com/DilemaFixer/gog/internal/search_engine"
)

func main() {
	enableLog := flag.Bool("log", false, "enable logging output")
	rootPath := flag.String("root", ".", "path to the working directory")
	flag.Parse()

	info, err := os.Stat(*rootPath)
	if os.IsNotExist(err) {
		fmt.Printf("Error: path '%s' does not exist\n", *rootPath)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("Error checking path: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Printf("Error: '%s' is not a directory\n", *rootPath)
		os.Exit(1)
	}

	var logger api.Logger
	if *enableLog {
		logger = log.NewChanLogger()
	} else {
		logger = log.NewNoopLogger()
	}

	eng := search_engine.NewSearchEngine(logger)
	err = eng.StartCyclicExecution(*rootPath)
	if err != nil {
		fmt.Printf("Error search engine: %v\n", err)
	}
}
